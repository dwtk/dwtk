package debugwire

func (dw *DebugWIRE) WriteSRAM(start uint16, b []byte) error {
	c := []byte{
		byte(start), byte(start >> 8),
	}
	if err := dw.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16((len(b) * 2) + 1)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x01,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x04,
		0x20,
	}
	return dw.device.Write(append(c, b...))
}

func (dw *DebugWIRE) ReadSRAM(start uint16, b []byte) error {
	c := []byte{
		byte(start), byte(start >> 8),
	}
	if err := dw.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16((len(b) * 2) + 1)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x00,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x00,
		0x20,
	}
	if err := dw.device.Write(c); err != nil {
		return err
	}

	return dw.device.Read(b)
}
