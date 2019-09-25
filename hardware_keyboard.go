package gemu 

const (
	CLEAR_BUFFER uint16 = 0
	GET_NEXT            = 1
	CHECK_KEY           = 2
	SET_INT             = 3
	SET_MODE            = 4
)

var keyboardClass = &HardwareClass{
	Name:  "keyboard",
	Desc:  "Generic Keyboard",
	DevID: 0x30c17406,
	VerID: 0x0001,
	MfgID: 0x1c6c8b36,
}

func init() {
	RegisterClass(keyboardClass)
}

type Keyboard struct {
	Hardware
	keycount  int
	keybuffer [8]uint8
	keydown   [8]uint8
	interrupt uint16
	mode      uint16
}

func NewKeyboard() *Keyboard {
	keyboard := &Keyboard{}
	keyboard.Class = keyboardClass
	return keyboard
}

func (K *Keyboard) HWI(D *DCPU) {
	switch D.Reg[0] {
	case CLEAR_BUFFER:
		K.keycount = 0
	case GET_NEXT:
		if K.keycount > 0 {
			D.Reg[2] = uint16(K.keybuffer[0])
			for i := 0; i < 7; i++ {
				K.keybuffer[i] = K.keybuffer[i+1]
			}
			K.keycount--
		} else {
			D.Reg[2] = 0
		}
	case CHECK_KEY:
		D.Reg[2] = 0
		for i := 0; i < 7; i++ {
			if K.keydown[i] == uint8(D.Reg[1]) {
				D.Reg[2] = 1
			}
		}
	case SET_INT:
		K.interrupt = D.Reg[1]
	case SET_MODE:
		K.mode = D.Reg[1]
		K.keycount = 0
		for i := range K.keydown {
			K.keydown[i] = 0
		}
	}
}

func (K *Keyboard) RawKey(key uint16, state bool) {
	if key >= 0x80 || K.mode == 1 {
		if state {
			handled := false
			for i := 0; i < 8; i++ {
				if K.keydown[i] == 0 {
					K.keydown[i] = uint8(key)
					handled = true
					break
				}
			}
			if !handled {
				K.keydown[7] = uint8(key)
			}
		} else {
			for i := 0; i < 8; i++ {
				if K.keydown[i] == uint8(key) {
					K.keydown[i] = 0
				}
			}
		}
		if !state {
			key |= 0x8000
		}
		if state || K.mode == 1 {
			K.queueKey(key)
		}
	}
}

func (K *Keyboard) ParsedKey(key uint16) {
	if K.mode == 0 {
		K.queueKey(key)
	}
}

func (K *Keyboard) queueKey(key uint16) {
	if K.keycount < 8 {
		K.keybuffer[K.keycount] = uint8(key)
		K.keycount++
	} else {
		copy(K.keybuffer[:], K.keybuffer[1:])
		K.keybuffer[7] = uint8(key)
	}
	if K.interrupt != 0 {
		if K.Up != nil {
			if dcpu, ok := K.Up.(*DCPU); ok {
				dcpu.Int(K.interrupt)
			}
		}
	}
}

func (K *Keyboard) Reset() {
	K.keycount = 0
	for i := range K.keydown {
		K.keydown[i] = 0
	}
}
