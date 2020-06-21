package debugwire

import (
	"github.com/dwtk/dwtk/avr"
)

func (dw *DebugWIRE) SetSP(b uint16) error {
	c := []byte{
		byte(b), byte(b >> 8),
	}
	return dw.WriteSRAM(avr.SPL, c)
}

func (dw *DebugWIRE) GetSP() (uint16, error) {
	c := make([]byte, 2)
	if err := dw.ReadSRAM(avr.SPL, c); err != nil {
		return 0, err
	}
	return (uint16(c[1]) << 8) | uint16(c[0]), nil
}

func (dw *DebugWIRE) SetSREG(b byte) error {
	return dw.WriteSRAM(avr.SREG, []byte{b})
}

func (dw *DebugWIRE) GetSREG() (byte, error) {
	c := make([]byte, 1)
	if err := dw.ReadSRAM(avr.SREG, c); err != nil {
		return 0, err
	}
	return c[0], nil
}
