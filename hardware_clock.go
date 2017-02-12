package gemu

import (
	"time"
)

var clockClass = &HardwareClass{
	Name:  "clock",
	Desc:  "Generic Clock",
	DevID: 0x12d0b402,
	VerID: 0x0002,
	MfgID: 0x1c6c8b36,
}

func init() {
	RegisterClass(clockClass)
}

type Clock struct {
	Hardware
	Rate         uint16
	RateAccum    uint16
	Total        uint16
	Accum        uint16
	TicksLeft    int
	Interrupt    uint16
	RealOffset   time.Duration
	RunTimeStart time.Time
}

func NewClock() *Clock {
	dev := &Clock{}
	dev.Class = clockClass
	realTime := time.Date(2600,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC)
	dev.RealOffset = realTime.Sub(time.Now())
	dev.RunTimeStart = time.Now()
	return dev
}

func (c *Clock) HWI(D *DCPU) {
	switch D.Reg[0] {
	case 0:
		c.Rate = D.Reg[1]
	case 1:
		D.Reg[2] = c.Total
		c.Total = 0
	case 2:
		c.Interrupt = D.Reg[1]
	case 0x0010:
		realTime := time.Now().Add(c.RealOffset)
		D.Reg[1] = uint16(realTime.Year())
		D.Reg[2] = uint16(int(realTime.Month())<<8 | realTime.Day())
		D.Reg[3] = uint16(realTime.Hour()<<8 | realTime.Minute())
		D.Reg[4] = uint16(realTime.Second())
		D.Reg[5] = uint16(realTime.Nanosecond() / int(time.Millisecond))
	case 0x0011:
		runTimeTotal := time.Now().Sub(c.RunTimeStart)
		D.Reg[2] = uint16(runTimeTotal / (time.Hour * 24))
		D.Reg[3] = uint16(runTimeTotal/(time.Hour))<<8 | uint16(runTimeTotal/(time.Minute))
		D.Reg[4] = uint16(runTimeTotal / (time.Second))
		D.Reg[5] = uint16(runTimeTotal / (time.Millisecond))
	case 0x0012:
		realTime := time.Date(int(D.Reg[1]),
			time.Month(D.Reg[2]>>8),
			int(D.Reg[2]&0xFF),
			int(D.Reg[3]>>8),
			int(D.Reg[3]&0xFF),
			int(D.Reg[4]),
			int(D.Reg[5])*int(time.Millisecond),
			time.UTC)
		c.RealOffset = realTime.Sub(time.Now())
	case 0xFFFF:
		c.Reset()
	}
}

func (c *Clock) Reset() {
	c.RunTimeStart = time.Now()
	c.Rate = 0
	c.Total = 0
	c.Interrupt = 0
}

func (c *Clock) Tick(ticks int) {
	c.TicksLeft -= ticks
	for c.TicksLeft < 0 {
		if c.Accum < 15 {
			c.Accum++
			c.TicksLeft += 1666
		} else {
			c.Accum = 0
			c.TicksLeft += 1676
		}
		c.RateAccum++
		if c.RateAccum >= c.Rate {
			c.RateAccum = 0
			if c.Rate > 0 {
				c.Total++
				if c.Interrupt != 0 && c.Up != nil {
					if dcpu, ok := c.Up.(*DCPU); ok {
						dcpu.Int(c.Interrupt)
					}
				}
			}
		}
	}
}
