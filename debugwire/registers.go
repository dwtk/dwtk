package debugwire

import (
	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWire) WriteRegisters(start byte, b []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(b)),
		0xc2, 0x05,
		0x20,
	}
	return dw.device.Write(append(c, b...))
}

func (dw *DebugWire) ReadRegisters(start byte, b []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(b)),
		0xc2, 0x01,
		0x20,
	}
	if err := dw.device.Write(c); err != nil {
		return err
	}
	return dw.device.Read(b)
}

func (dw *DebugWire) SetPC(b uint16) error {
	dw.afterBreak = false
	b /= 2
	c := []byte{
		0xd0, byte(b >> 8), byte(b),
	}
	return dw.device.Write(c)
}

func (dw *DebugWire) GetPC() (uint16, error) {
	if err := dw.device.Write([]byte{0xf0}); err != nil {
		return 0, err
	}

	rv, err := dw.device.ReadWord()
	if err != nil {
		return 0, err
	}

	if dw.afterBreak {
		if rv > 0 {
			rv -= 1
		}
		dw.afterBreak = false
	}

	rv *= 2
	return rv, nil
}

func (dw *DebugWire) SetSP(b uint16) error {
	c := []byte{
		byte(b), byte(b >> 8),
	}
	return dw.WriteSRAM(avr.SPL, c)
}

func (dw *DebugWire) GetSP() (uint16, error) {
	c := make([]byte, 2)
	if err := dw.ReadSRAM(avr.SPL, c); err != nil {
		return 0, err
	}
	return uint16(c[1]<<8) | uint16(c[0]), nil
}

func (dw *DebugWire) SetSREG(b byte) error {
	return dw.WriteSRAM(avr.SREG, []byte{b})
}

func (dw *DebugWire) GetSREG() (byte, error) {
	c := make([]byte, 1)
	if err := dw.ReadSRAM(avr.SREG, c); err != nil {
		return 0, err
	}
	return c[0], nil
}
