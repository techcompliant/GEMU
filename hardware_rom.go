package gemu

import (
	"reflect"
	"unsafe"
)

var romClass = &HardwareClass{
	Name:  "rom",
	Desc:  "Embedded ROM",
	DevID: 0x17400011,
	VerID: 0x0001,
	MfgID: 0x12452135,
}

func init() {
	RegisterClass(romClass)
}

type ROM struct {
	Hardware
	Data []uint16
}

func NewRom(romImage string, flip bool) *ROM {
	rom := &ROM{}
	rom.Class = romClass
	storage := defaultStorage
	if flip {
		storage = NewFlipStorage(storage)
	}
	if storage.Exists(romImage) {
		rom.Data = make([]uint16, storage.Length(romImage)/2)
		rawData := []byte{}
		bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&rawData))
		bytesHeader.Data = uintptr(unsafe.Pointer(&rom.Data[0]))
		bytesHeader.Len = len(rom.Data) * 2
		bytesHeader.Cap = len(rom.Data) * 2
		storage.Read(romImage, 0, rawData)
	}
	return rom
}

func (R *ROM) HWI(D *DCPU) {
	switch D.Reg[0] {
	case 0:
		D.Mem.LoadMem(R.Data)
		D.PC = 0
	case 1:
		copy(R.Data, D.Mem.GetRaw()[D.Reg[1]:])
	}
}

func (R *ROM) Reset() {
	if R.Up != nil {
		if dcpu, ok := R.Up.(*DCPU); ok {
			dcpu.Mem.LoadMem(R.Data)
		}
	}
}
