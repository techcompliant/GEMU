package gemu // import "gemu.techcompliant.com/gemu"

type Hardware struct {
	Up    IHardware
	Down  []IHardware
	Class *HardwareClass
}

func (H *Hardware) GetClass() *HardwareClass {
	return H.Class
}

func (H *Hardware) GetUp() IHardware {
	return H.Up
}

func (H *Hardware) GetMem() IMem {
	if H.Up != nil {
		if dcpu, ok := H.Up.(*DCPU); ok {
			return dcpu.Mem
		}
	}
	return nil
}

func (H *Hardware) SetUp(up IHardware) {
	H.Up = up
}

func (H *Hardware) Attach(Add IHardware) {
	H.Down = append(H.Down, Add)
}

func (H *Hardware) GetDown() []IHardware {
	return H.Down
}

func (H *Hardware) Start() {
	for _, obj := range H.Down {
		obj.Start()
	}
}

func (H *Hardware) Stop() {
	for _, obj := range H.Down {
		obj.Stop()
	}
}

func (H *Hardware) Reset() {
	for _, obj := range H.Down {
		obj.Reset()
	}
}

func (H *Hardware) HWI(D *DCPU) {
	//log.Printf("Unhandled HWI\n")
}

func (H *Hardware) HWQ(D *DCPU) {
	C := H.GetClass()
	D.Reg[0] = uint16(C.DevID & 0xFFFF)
	D.Reg[1] = uint16((C.DevID >> 16) & 0xFFFF)
	D.Reg[2] = C.VerID
	D.Reg[3] = uint16(C.MfgID & 0xFFFF)
	D.Reg[4] = uint16((C.MfgID >> 16) & 0xFFFF)
}

type IHardware interface {
	GetUp() IHardware
	GetMem() IMem
	Attach(Add IHardware)
	GetDown() []IHardware
	SetUp(up IHardware)
	GetClass() *HardwareClass
	Start()
	Stop()
	Reset()
	HWI(D *DCPU)
	HWQ(D *DCPU)
}

type Ticker interface {
	Tick(int)
}

type HardwareClass struct {
	Name  string
	Desc  string
	DevID uint32
	VerID uint16
	MfgID uint32
}

type IStateChanges interface {
	IsDirty() bool
	ClearDirty()
}

var Classes []*HardwareClass

func RegisterClass(hc *HardwareClass) {
	Classes = append(Classes, hc)
}

type IMem interface {
	ReadMem(addr uint16) uint16
	WriteMem(addr uint16, val uint16)
	LoadMem(data []uint16)
	GetRaw() []uint16
	RegisterSync(addr uint16, synclen uint16) *Sync
	Reset()
}
