package gemu

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Instruction struct {
	Opcode   uint16
	OpA      uint16
	AddA     uint16
	OpB      uint16
	AddB     uint16
	SPBShift bool
}

func OperandString(op uint16, add uint16) string {
	switch op {
	case 0x00:
		return "A"
	case 0x01:
		return "B"
	case 0x02:
		return "C"
	case 0x03:
		return "X"
	case 0x04:
		return "Y"
	case 0x05:
		return "Z"
	case 0x06:
		return "I"
	case 0x07:
		return "J"
	case 0x08:
		return "[A]"
	case 0x09:
		return "[B]"
	case 0x0A:
		return "[C]"
	case 0x0B:
		return "[X]"
	case 0x0C:
		return "[Y]"
	case 0x0D:
		return "[Z]"
	case 0x0E:
		return "[I]"
	case 0x0F:
		return "[J]"
	case 0x10:
		return fmt.Sprintf("[A+%04x]", add)
	case 0x11:
		return fmt.Sprintf("[B+%04x]", add)
	case 0x12:
		return fmt.Sprintf("[C+%04x]", add)
	case 0x13:
		return fmt.Sprintf("[X+%04x]", add)
	case 0x14:
		return fmt.Sprintf("[Y+%04x]", add)
	case 0x15:
		return fmt.Sprintf("[Z+%04x]", add)
	case 0x16:
		return fmt.Sprintf("[I+%04x]", add)
	case 0x17:
		return fmt.Sprintf("[J+%04x]", add)
	case 0x18:
		return "PSHPOP"
	case 0x19:
		return "[SP]"
	case 0x1a:
		return fmt.Sprintf("[SP+%04x]", add)
	case 0x1b:
		return "SP"
	case 0x1C:
		return "PC"
	case 0x1D:
		return "EX"
	case 0x1E:
		return fmt.Sprintf("[%04x]", add)
	case 0x1F:
		return fmt.Sprintf("%04x", add)
	default:
		return fmt.Sprintf("%04x", op-0x21)
	}
}

func (I *Instruction) String() string {
	Inst := fmt.Sprintf("ERR:%02x", I.Opcode)
	if I.Opcode != 0 {
		switch I.Opcode {
		case 0x01:
			Inst = "SET"
		case 0x02:
			Inst = "ADD"
		case 0x03:
			Inst = "SUB"
		case 0x04:
			Inst = "MUL"
		case 0x05:
			Inst = "MLI"
		case 0x06:
			Inst = "DIV"
		case 0x07:
			Inst = "DVI"
		case 0x08:
			Inst = "MOD"
		case 0x09:
			Inst = "MDI"
		case 0x0A:
			Inst = "AND"
		case 0x0B:
			Inst = "BOR"
		case 0x0C:
			Inst = "XOR"
		case 0x0D:
			Inst = "SHR"
		case 0x0E:
			Inst = "ASR"
		case 0x0F:
			Inst = "SHL"
		case 0x10:
			Inst = "IFB"
		case 0x11:
			Inst = "IFC"
		case 0x12:
			Inst = "IFE"
		case 0x13:
			Inst = "IFN"
		case 0x14:
			Inst = "IFG"
		case 0x15:
			Inst = "IFA"
		case 0x16:
			Inst = "IFL"
		case 0x17:
			Inst = "IFU"
		case 0x1A:
			Inst = "ADX"
		case 0x1B:
			Inst = "SBX"
		case 0x1E:
			Inst = "STI"
		case 0x1F:
			Inst = "STD"
		}
		Inst += " " + OperandString(I.OpB, I.AddB) + " " + OperandString(I.OpA, I.AddA)
	} else {
		Inst = fmt.Sprintf("SPC:%02x", I.OpB)
		switch I.OpB {
		case 0x01:
			Inst = "JSR"
		case 0x08:
			Inst = "INT"
		case 0x09:
			Inst = "IAG"
		case 0x0A:
			Inst = "IAS"
		case 0x0B:
			Inst = "RFI"
		case 0x0C:
			Inst = "IAQ"
		case 0x10:
			Inst = "HWN"
		case 0x11:
			Inst = "HWQ"
		case 0x12:
			Inst = "HWI"
		case 0x13:
			Inst = "LOG"
		case 0x14:
			Inst = "BRK"
		case 0x15:
			Inst = "HLT"
		}
		Inst += " " + OperandString(I.OpA, I.AddA)
	}
	return Inst
}

func Decode(inst uint16) Instruction {
	ret := Instruction{}
	ret.Opcode = inst & 0x001F
	ret.OpA = inst >> 10
	ret.OpB = (inst >> 5) & 0x001F
	return ret
}

func (I *Instruction) PreProcessOps(D *DCPU) {
	switch {
	case I.OpA >= 0x10 && I.OpA <= 0x17:
		fallthrough
	case I.OpA == 0x1a || I.OpA == 0x1e || I.OpA == 0x1f:
		I.AddA = D.Mem.ReadMem(D.PC)
		D.PC++
	}
	if I.Opcode != 0 {
		switch {
		case I.OpB >= 0x10 && I.OpB <= 0x17:
			fallthrough
		case I.OpB == 0x1a || I.OpB == 0x1e || I.OpB == 0x1f:
			I.AddB = D.Mem.ReadMem(D.PC)
			D.PC++
		}
	}
}

func (I *Instruction) GetOp(D *DCPU, OpB bool) (val uint16) {
	Op := I.OpA
	Add := I.AddA
	if OpB {
		Op = I.OpB
		Add = I.AddB
	}
	switch {
	case Op <= 0x07:
		return D.Reg[Op]
	case Op <= 0x0F:
		return D.Mem.ReadMem(D.Reg[Op-0x08])
	case Op <= 0x17:
		D.WaitState++
		return D.Mem.ReadMem(D.Reg[Op-0x10] + Add)
	case Op == 0x18:
		if OpB {
			I.SPBShift = true
			D.SP--
			return D.Mem.ReadMem(D.SP)
		} else {
			D.SP++
			return D.Mem.ReadMem(D.SP - 1)
		}
	case Op == 0x19:
		return D.Mem.ReadMem(D.SP)
	case Op <= 0x1A:
		D.WaitState++
		return D.Mem.ReadMem(D.SP + Add)
	case Op == 0x1B:
		return D.SP
	case Op == 0x1C:
		return D.PC
	case Op == 0x1D:
		return D.EX
	case Op == 0x1E:
		D.WaitState++
		return D.Mem.ReadMem(Add)
	case Op == 0x1F:
		D.WaitState++
		return Add
	case Op <= 0x3F:
		return Op - 0x21
	}
	panic(fmt.Sprintf("Invalid operand: %s", I.String()))
}

func (I *Instruction) SetOp(D *DCPU, val uint16) {
	Op := I.OpB
	Add := I.AddB
	switch {
	case Op <= 0x07:
		D.Reg[Op] = val
		return
	case Op <= 0x0F:
		D.Mem.WriteMem(D.Reg[Op-0x08], val)
		return
	case Op <= 0x17:
		D.Mem.WriteMem(D.Reg[Op-0x10]+Add, val)
		D.WaitState++
		return
	case Op == 0x18:
		if !I.SPBShift {
			D.SP--
		}
		D.Mem.WriteMem(D.SP, val)
		return
	case Op == 0x19:
		D.Mem.WriteMem(D.SP, val)
		return
	case Op <= 0x1A:
		D.Mem.WriteMem(D.SP+Add, val)
		D.WaitState++
		return
	case Op == 0x1B:
		D.SP = val
		return
	case Op == 0x1C:
		D.PC = val
		return
	case Op == 0x1D:
		D.EX = val
		return
	case Op == 0x1E:
		D.Mem.WriteMem(Add, val)
		D.WaitState++
		return
	case Op == 0x1f:
		return
	}

	panic(fmt.Sprintf("Invalid operand: %s", I.String()))
}

func (I *Instruction) Run(D *DCPU) {
	I.PreProcessOps(D)
	if D.Skipping {
		if I.Opcode >= 0x10 && I.Opcode <= 0x17 {
			return
		}
		D.Skipping = false
		return
	}
	switch I.Opcode {
	case 0x01: // SET
		I.SetOp(D, I.GetOp(D, false))
	case 0x02: // ADD
		D.WaitState = 1
		a, b := uint32(I.GetOp(D, false)), uint32(I.GetOp(D, true))
		val := b + a
		I.SetOp(D, uint16(val&0xFFFF))
		D.EX = 0
		if val > 0xFFFF {
			D.EX = 1
		}
	case 0x03: // SUB
		D.WaitState = 1
		a, b := int32(I.GetOp(D, false)), int32(I.GetOp(D, true))
		val := b - a
		I.SetOp(D, uint16(val))
		D.EX = 0
		if val < 0 {
			D.EX = 0xFFFF
		}
	case 0x04: // MUL
		D.WaitState = 1
		a, b := uint32(I.GetOp(D, false)), uint32(I.GetOp(D, true))
		val := b * a
		I.SetOp(D, (uint16)(val&0xFFFF))
		D.EX = uint16((val >> 16) & 0xFFFF)
	case 0x05: // MLI
		D.WaitState = 1
		a, b := int32(int16(I.GetOp(D, false))), int32(int16(I.GetOp(D, true)))
		val := b * a
		I.SetOp(D, (uint16)(val&0xFFFF))
		D.EX = uint16((val >> 16) & 0xFFFF)
	case 0x06: // DIV
		D.WaitState = 2
		a, b := uint32(I.GetOp(D, false)), uint32(I.GetOp(D, true))
		if a == 0 {
			I.SetOp(D, 0)
			D.EX = 0
		} else {
			val := (b << 16) / a
			I.SetOp(D, (uint16)((val>>16)&0xFFFF))
			D.EX = (uint16)(val & 0xFFFF)
		}
	case 0x07: // DVI
		D.WaitState = 2
		a, b := int32(int16(I.GetOp(D, false))), int32(int16(I.GetOp(D, true)))
		if a == 0 {
			I.SetOp(D, 0)
			D.EX = 0
		} else {
			sign := (a < 0) != (b < 0)
			if a < 0 {
				a = -a
			}
			if b < 0 {
				b = -b
			}
			val := (b << 16) / a
			if sign {
				val = -val
			}
			I.SetOp(D, (uint16)((val>>16)&0xFFFF))
			D.EX = (uint16)(val & 0xFFFF)
		}
	case 0x08: // MOD
		D.WaitState = 2
		a, b := uint16(I.GetOp(D, false)), uint16(I.GetOp(D, true))
		if a == 0 {
			b = 0
		} else {
			b = b % a
		}
		I.SetOp(D, uint16(b))
	case 0x09: // MDI
		D.WaitState = 2
		a, b := int16(I.GetOp(D, false)), int16(I.GetOp(D, true))
		if a == 0 {
			b = 0
		} else {
			b = b % a
		}
		I.SetOp(D, uint16(b))
	case 0x0A: // AND
		val := I.GetOp(D, false) & I.GetOp(D, true)
		I.SetOp(D, val)
	case 0x0B: // BOR
		val := I.GetOp(D, false) | I.GetOp(D, true)
		I.SetOp(D, val)
	case 0x0C: // XOR
		val := I.GetOp(D, false) ^ I.GetOp(D, true)
		I.SetOp(D, val)
	case 0x0D: // SHR
		a, b := I.GetOp(D, false), I.GetOp(D, true)
		a = a & 31
		I.SetOp(D, b>>a)
		D.EX = uint16(((uint32(b) << 16) >> a) & 0xFFFF)
	case 0x0E: // ASR
		a, b := I.GetOp(D, false), int16(I.GetOp(D, true))
		a = a & 31
		I.SetOp(D, uint16(b>>a))
		D.EX = uint16((int32(uint32(b)<<16) >> a) & 0xFFFF)
	case 0x0F: // SHL
		a, b := I.GetOp(D, false), I.GetOp(D, true)
		a = a & 31
		I.SetOp(D, b<<a)
		D.EX = uint16(((uint32(b) << a) >> 16) & 0xFFFF)
	case 0x10: // IFB
		D.WaitState = 1
		D.Skipping = (I.GetOp(D, false) & I.GetOp(D, true)) == 0
	case 0x11: // IFC
		D.WaitState = 1
		D.Skipping = (I.GetOp(D, false) & I.GetOp(D, true)) != 0
	case 0x12: // IFE
		D.WaitState = 1
		D.Skipping = I.GetOp(D, false) != I.GetOp(D, true)
	case 0x13: // IFN
		D.WaitState = 1
		D.Skipping = I.GetOp(D, false) == I.GetOp(D, true)
	case 0x14: // IFG
		D.WaitState = 1
		D.Skipping = I.GetOp(D, false) >= I.GetOp(D, true)
	case 0x15: // IFA
		D.WaitState = 1
		D.Skipping = int16(I.GetOp(D, false)) >= int16(I.GetOp(D, true))
	case 0x16: // IFL
		D.WaitState = 1
		D.Skipping = I.GetOp(D, false) <= I.GetOp(D, true)
	case 0x17: // IFU
		D.WaitState = 1
		D.Skipping = int16(I.GetOp(D, false)) <= int16(I.GetOp(D, true))
	case 0x1a: // ADX
		D.WaitState = 2
		a, b := uint32(I.GetOp(D, false)), uint32(I.GetOp(D, true))
		val := a + b + uint32(D.EX)
		I.SetOp(D, uint16(val&0xFFFF))
		D.EX = 0
		if val > 0xFFFF {
			D.EX = 1
		}
	case 0x1b: // SBX
		D.WaitState = 2
		a, b := uint32(I.GetOp(D, false)), uint32(I.GetOp(D, true))
		val := b - a + uint32(D.EX)
		I.SetOp(D, uint16(val))
		D.EX = 0
		if val < 0 {
			D.EX = 0xFFFF
		}
	case 0x1E: // STI
		D.WaitState = 1
		I.SetOp(D, I.GetOp(D, false))
		D.Reg[6]++
		D.Reg[7]++
	case 0x1F: // STD
		D.WaitState = 1
		I.SetOp(D, I.GetOp(D, false))
		D.Reg[6]--
		D.Reg[7]--
	case 0x00: // SPECIAL OPCODES
		switch I.OpB {
		case 0x01: // JSR
			D.WaitState = 2
			dest := I.GetOp(D, false)
			D.SP--
			D.Mem.WriteMem(D.SP, D.PC)
			D.PC = dest
		case 0x08: // INT
			D.WaitState = 3
			D.Int(I.GetOp(D, false))
		case 0x09: // IAG
			orig := I.OpB
			I.OpB = I.OpA
			I.AddB = I.AddA
			I.SetOp(D, D.IA)
			I.OpB = orig
		case 0x0A: // IAS
			D.WaitState = 2
			D.IA = I.GetOp(D, false)
		case 0x0B: // RFI
			D.WaitState = 1
			D.EnIQ = true
			D.Reg[0] = D.Mem.ReadMem(D.SP)
			D.PC = D.Mem.ReadMem(D.SP + 1)
			D.SP += 2
		case 0x0C: // IAQ
			D.EnIQ = I.GetOp(D, false) == 0
		case 0x10: // HWN
			D.WaitState = 1
			orig := I.OpB
			I.OpB = I.OpA
			I.AddB = I.AddA
			I.SetOp(D, uint16(len(D.Down)))
			I.OpB = orig
		case 0x11: // HWQ
			D.WaitState = 3
			id := int(I.GetOp(D, false))
			D.Reg[0] = 0
			D.Reg[1] = 0
			D.Reg[2] = 0
			D.Reg[3] = 0
			D.Reg[4] = 0
			if id < len(D.Down) && D.Down[id] != nil {
				D.Down[id].HWQ(D)
			}
		case 0x12: // HWI
			D.WaitState = 3
			id := int(I.GetOp(D, false))
			if id < len(D.Down) && D.Down[id] != nil {
				D.Down[id].HWI(D)
			}
		case 0x13: // LOG
			//log.Printf("DCPU Log: %04x\n", I.GetOp(D, false))
		case 0x14: // BRK
		case 0x15: // HLT
			D.WaitInt = true
		default:
			//log.Printf("DCPU: Invalid opcode: %04x\n", I.Opcode)
			D.Running = false
		}
	default:
		//log.Printf("DCPU: Invalid opcode: %04x\n", I.Opcode)
		D.Running = false
	}
}

type DCPU struct {
	Hardware
	Reg        [8]uint16
	PC         uint16
	SP         uint16
	EX         uint16
	IA         uint16
	IQLen      uint16
	IQ         [256]uint16
	EnIQ       bool
	WaitState  int
	Skipping   bool
	Running    bool
	WaitInt    bool
	TickRate   int
	SpareTicks int

	Mem *Mem16x64k
}

var dcpuClass = &HardwareClass{
	Name: "dcpu",
	Desc: "DCPU-16 1.7",
}

func init() {
	RegisterClass(dcpuClass)
}

func NewDCPU(cycleRate int) *DCPU {
	if cycleRate == 0 {
		cycleRate = 100000
	}
	tickRate := 1
	if cycleRate > 100000 {
		//log.Printf("Error: cycle rate above 100Khz requested!\n")
	} else {
		tickRate = int(100000 / cycleRate)
	}

	dcpu := &DCPU{TickRate: tickRate, Mem: &Mem16x64k{}}

	dcpu.Mem.RawRAM = []byte{}
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&dcpu.Mem.RawRAM))
	bytesHeader.Data = uintptr(unsafe.Pointer(&dcpu.Mem.RAM[0]))
	bytesHeader.Len = 65536 * 2
	bytesHeader.Cap = 65536 * 2

	dcpu.Class = dcpuClass

	return dcpu
}

func (D *DCPU) String() string {
	return fmt.Sprintf("A: %04x B: %04x C: %04x X: %04x Y: %04x Z: %04x I: %04x J: %04x PC: %04x SP: %04x EX: %04x IA: %04x",
		D.Reg[0], D.Reg[1], D.Reg[2], D.Reg[3], D.Reg[4], D.Reg[5], D.Reg[6], D.Reg[7],
		D.PC, D.SP, D.EX, D.IA)
}

func (D *DCPU) Tick(ticks int) {
	if !D.Running {
		return
	}
	ticks, D.SpareTicks = (ticks+D.SpareTicks)/D.TickRate, (ticks+D.SpareTicks)%D.TickRate

	for l1 := 0; l1 < ticks; l1++ {
		//DCPUTick++
		if D.WaitState > 0 {
			D.WaitState--
			continue
		}
		if !D.Skipping && D.EnIQ && D.IQLen > 0 {
			D.IQLen--
			D.EnIQ = false
			D.Mem.WriteMem(D.SP-1, D.PC)
			D.Mem.WriteMem(D.SP-2, D.Reg[0])
			D.SP -= 2
			D.PC = D.IA
			D.Reg[0] = D.IQ[0]
			copy(D.IQ[:], D.IQ[1:])
			D.WaitInt = false
		}
		if D.WaitInt || !D.Running {
			return
		}

		I := Decode(D.Mem.ReadMem(D.PC))

		D.PC++
		I.Run(D)

	}
}

func (D *DCPU) Start() {
	D.Reset()
	D.Hardware.Start()
	D.Running = true
}

func (D *DCPU) Stop() {
	D.Running = false
	D.Hardware.Stop()
}

func (D *DCPU) Reset() {
	for i := range D.Reg {
		D.Reg[i] = 0
	}
	D.PC = 0
	D.IA = 0
	D.EX = 0
	D.SP = 0
	D.IQLen = 0
	D.EnIQ = true
	if D.Mem != nil {
		D.Mem.Reset()
	}
	D.Hardware.Reset()
}

func (D *DCPU) Int(msg uint16) {
	if D.IA == 0 {
		return
	}
	if D.IQLen >= 256 {
		return
	}
	D.IQ[D.IQLen] = msg
	D.IQLen++
}

func (D *DCPU) GetMem() IMem {
	return D.Mem
}

type Mem16x64k struct {
	RAM       [65536]uint16
	RawRAM    []byte
	Dirty     [4096]uint8
	SyncCount [4096]uint8
}

type Sync struct {
	addr       uint16
	synclen    uint16
	registered bool
	mem        *Mem16x64k
}

func (S *Sync) Register(mem *Mem16x64k) {
	if !S.registered {
		for l1 := S.addr >> 4; l1 <= ((S.addr + S.synclen) >> 4); l1++ {
			mem.SyncCount[l1]++
			mem.Dirty[l1] = 1
		}
		S.registered = true
		S.mem = mem
	}
}

func (S *Sync) Unregister() {
	if S.registered {
		for l1 := S.addr >> 4; l1 <= ((S.addr + S.synclen) >> 4); l1++ {
			S.mem.SyncCount[l1]--
		}
		S.registered = false
		S.mem = nil
	}
}

func (M *Mem16x64k) ReadMem(addr uint16) uint16 {
	if int(addr) > cap(M.RAM) {
		//log.Fatal(errors.New("Out of bounds memory  ReadMem"))
	}
	return M.RAM[addr]
}

func (M *Mem16x64k) WriteMem(addr uint16, val uint16) {
	if int(addr) > cap(M.RAM) {
		//log.Fatal(errors.New("Out of bounds memory WriteMem"))
	}
	M.Dirty[addr>>4] = 1
	M.RAM[addr] = val
}

func (M *Mem16x64k) LoadMem(data []uint16) {
	copy(M.RAM[:], data)
}

func (M *Mem16x64k) GetRaw() []uint16 {
	return M.RAM[:]
}

func (M *Mem16x64k) RegisterSync(addr uint16, synclen uint16) *Sync {
	ret := &Sync{
		addr:    addr,
		synclen: synclen,
	}
	ret.Register(M)
	return ret
}

func (M *Mem16x64k) Reset() {
	for i := range M.RAM {
		M.RAM[i] = 0
	}
}
