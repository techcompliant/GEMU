package gemu

type EEPROM struct {
	Data []uint16
}

func (E *EEPROM) HWI(D *DCPU) {
	switch D.Reg[1] {
	case 1:
		D.Reg[4] = E.Data[D.Reg[3]]
		D.WaitState += 1000 / D.TickRate
	case 2:
		E.Data[D.Reg[3]] &= D.Reg[4]
		D.WaitState += 5000 / D.TickRate
	case 3:
		for i := range E.Data {
			E.Data[i] = 0xFFFF
		}
		D.WaitState += 10000 / D.TickRate
	}
}
