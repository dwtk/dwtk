package debugwire

func (dw *DebugWire) WriteRegisters(start byte, b []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(b)),
		0xc2, 0x05,
		0x20,
	}
	return dw.Port.Write(append(c, b...))
}

func (dw *DebugWire) ReadRegisters(start byte, b []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(b)),
		0xc2, 0x01,
		0x20,
	}
	if err := dw.Port.Write(c); err != nil {
		return err
	}
	return dw.Port.Read(b)
}

func (dw *DebugWire) SetPC(b uint16) error {
	dw.afterBreak = false
	b /= 2
	c := []byte{
		0xd0, byte(b >> 8), byte(b),
	}
	return dw.Port.Write(c)
}

func (dw *DebugWire) GetPC() (uint16, error) {
	if err := dw.Port.Write([]byte{0xf0}); err != nil {
		return 0, err
	}

	rv, err := dw.Port.ReadWord()
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
