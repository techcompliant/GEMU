package gemu

var pixieClass = &HardwareClass{
	Name:  "pixie",
	Desc:  "PIXIE",
	DevID: 0x734df615,
	VerID: 0x1802,
	MfgID: 0x1c6c8b36,
}

func init() {
	RegisterClass(pixieClass)
}

func NewPIXIE() *PIXIE {
	pixie := &PIXIE{}
	pixie.Class = pixieClass
	pixie.NeedSync = true
	return pixie
}

func (P *PIXIE) HWI(D *DCPU) {
	switch D.Reg[0] {
	case 0:
		if P.dspSync != nil {
			P.dspSync.Unregister()
			P.dspSync = nil
		}
		P.DspMem = D.Reg[1]
		if P.DspMem != 0 {
			syncCount := uint16(384)
			if P.Mode >= 1 && P.Mode <= 4 {
				syncCount = P.Mode * 768
			}
			P.dspSync = D.Mem.RegisterSync(D.Reg[1], syncCount)
		}
		P.NeedSync = true
	case 1:
		if P.fontSync != nil {
			P.fontSync.Unregister()
			P.fontSync = nil
		}
		P.FontMem = D.Reg[1]
		if P.FontMem != 0 {
			P.fontSync = D.Mem.RegisterSync(D.Reg[1], 256)
		}
		P.NeedSync = true
	case 2:
		if P.palSync != nil {
			P.palSync.Unregister()
			P.palSync = nil
		}
		P.PalMem = D.Reg[1]
		if P.PalMem != 0 {
			P.palSync = D.Mem.RegisterSync(D.Reg[1], 16)
		}
		P.NeedSync = true
	case 3:
		P.Border = D.Reg[1]
		P.NeedSync = true
	case 4:
		copy(P.GetMem().GetRaw()[D.Reg[1]:], LemDefFont)
	case 5:
		copy(P.GetMem().GetRaw()[D.Reg[1]:], LemDefPal)
	case 6:
		P.Mode = D.Reg[1]
		if P.DspMem != 0 {
			P.dspSync.Unregister()
			syncCount := uint16(384)
			if P.Mode >= 1 && P.Mode <= 4 {
				syncCount = P.Mode * 768
			}
			P.dspSync = D.Mem.RegisterSync(P.DspMem, syncCount)
		}
		if P.Mode > 0 {
			if P.fontSync != nil {
				P.fontSync.Unregister()
				P.fontSync = nil
			}
		} else {
			if P.fontSync == nil && P.FontMem != 0 {
				P.fontSync = D.Mem.RegisterSync(P.FontMem, 256)
			}
		}
	}
}

type PIXIE struct {
	Hardware
	NeedSync bool

	DspMem  uint16
	FontMem uint16
	PalMem  uint16
	Border  uint16
	Mode    uint16

	dspSync  *Sync
	fontSync *Sync
	palSync  *Sync
}

func (P *PIXIE) Reset() {
	P.DspMem = 0
	P.FontMem = 0
	P.PalMem = 0
	P.Border = 0
	P.Mode = 0
	if P.dspSync != nil {
		P.dspSync.Unregister()
	}
	P.dspSync = nil
	if P.fontSync != nil {
		P.fontSync.Unregister()
	}
	P.fontSync = nil
	if P.palSync != nil {
		P.palSync.Unregister()
	}
	P.palSync = nil
	P.NeedSync = true
}

func (P *PIXIE) IsDirty() bool {
	return P.NeedSync
}

func (P *PIXIE) ClearDirty() {
	P.NeedSync = false
}
