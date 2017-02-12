package gemu

import (
	"reflect"
	"unsafe"
)

var floppyClass = &HardwareClass{
	Name:  "mack_35fd",
	Desc:  "Mackapar M35FD",
	DevID: 0x4fd524c5,
	VerID: 0x000b,
	MfgID: 0x1eb37e91,
}

func init() {
	RegisterClass(floppyClass)
}

type M35FD struct {
	Hardware
	Error     uint16
	interrupt uint16

	Running    bool
	TicksLeft  int
	Block      []byte
	Addr       uint16
	Read       bool
	ActionDone bool

	Disk string

	NeedSync bool

	storage Storage
}

func NewM35FD(flip bool) *M35FD {
	floppy := &M35FD{}
	floppy.Class = floppyClass
	floppy.NeedSync = true
	floppy.storage = defaultStorage
	if flip {
		floppy.storage = NewFlipStorage(defaultStorage)
	}
	return floppy
}

func (fd *M35FD) HWI(D *DCPU) {
	switch D.Reg[0] {
	case 0:
		if fd.Running {
			D.Reg[1] = 3
		} else if fd.Disk != "" {
			/*if disk.Readonly {
				D.Reg[1] = 2
			} else {
				D.Reg[1] = 1
			}*/
			D.Reg[1] = 1 // Need to reenable readonly checks
		} else {
			D.Reg[1] = 0
		}
		D.Reg[2] = fd.Error
		fd.Error = 0
	case 1:
		fd.interrupt = D.Reg[3]
	case 2:
		if fd.Disk != "" {
			fd.Block = make([]byte, 1024)
			fd.ActionDone = false
			go func(offset int) { fd.storage.Read(fd.Disk, offset, fd.Block); fd.ActionDone = true }(int(D.Reg[3]) * 1024)
			fd.Running = true
			fd.TicksLeft = 10000
			fd.Addr = D.Reg[4]
			fd.Read = true
			fd.NeedSync = true
			D.Reg[1] = 1
		} else {
			fd.Error = 2
			D.Reg[1] = 0
		}
		if fd.interrupt != 0 && fd.Up != nil {
			if dcpu, ok := fd.Up.(*DCPU); ok {
				dcpu.Int(fd.interrupt)
			}
		}
	case 3:
		if fd.Disk != "" {

			fd.Block = make([]byte, 1024)
			rawData := []uint16{}
			bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&rawData))
			bytesHeader.Data = uintptr(unsafe.Pointer(&fd.Block[0]))
			bytesHeader.Len = 512
			bytesHeader.Cap = 512
			baseRam := fd.GetMem().GetRaw()[fd.Addr:]
			copy(rawData, baseRam)
			go func(offset int) { fd.storage.Write(fd.Disk, offset, fd.Block); fd.ActionDone = true }(int(D.Reg[3]) * 1024)
			fd.Running = true
			fd.TicksLeft = 10000
			fd.Addr = D.Reg[4]
			fd.Read = false
			fd.NeedSync = true
			D.Reg[1] = 1
		} else {
			fd.Error = 2
			D.Reg[1] = 0
		}
		if fd.interrupt != 0 && fd.Up != nil {
			if dcpu, ok := fd.Up.(*DCPU); ok {
				dcpu.Int(fd.interrupt)
			}
		}
	}
}

func (fd *M35FD) Tick(ticks int) {
	if fd.Running {
		fd.TicksLeft -= ticks
		if fd.TicksLeft <= 0 {
			if fd.Read {
				if fd.ActionDone {
					rawData := []uint16{}
					bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&rawData))
					bytesHeader.Data = uintptr(unsafe.Pointer(&fd.Block[0]))
					bytesHeader.Len = 512
					bytesHeader.Cap = 512
					if fd.GetMem() != nil {
						copy(fd.GetMem().GetRaw()[fd.Addr:], rawData)
					}
					fd.Running = false
					fd.NeedSync = true
					if fd.interrupt != 0 && fd.Up != nil {
						if dcpu, ok := fd.Up.(*DCPU); ok {
							dcpu.Int(fd.interrupt)
						}
					}
				}
			} else {
				if fd.ActionDone {
					fd.Running = false
					fd.NeedSync = true
					if fd.interrupt != 0 && fd.Up != nil {
						if dcpu, ok := fd.Up.(*DCPU); ok {
							dcpu.Int(fd.interrupt)
						}
					}
				}
			}

		}
	}
}

func (fd *M35FD) Reset() {
	fd.Error = 0
}

func (fd *M35FD) ChangeDisk(disk string) {
	fd.Disk = disk
	if fd.interrupt != 0 && fd.Up != nil {
		if dcpu, ok := fd.Up.(*DCPU); ok {
			dcpu.Int(fd.interrupt)
		}
	}
}

func (fd *M35FD) IsDirty() bool {
	return fd.NeedSync
}

func (fd *M35FD) ClearDirty() {
	fd.NeedSync = false
}
