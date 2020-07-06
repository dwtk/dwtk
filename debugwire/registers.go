package debugwire

func (dw *DebugWIRE) SetSP(b uint16) error {
	c := []byte{
		byte(b), byte(b >> 8),
	}
	return dw.WriteSRAM(dw.MCU.SP().Mem16(), c)
}

func (dw *DebugWIRE) GetSP() (uint16, error) {
	c := make([]byte, 2)
	if err := dw.ReadSRAM(dw.MCU.SP().Mem16(), c); err != nil {
		return 0, err
	}
	return (uint16(c[1]) << 8) | uint16(c[0]), nil
}

func (dw *DebugWIRE) SetSREG(b byte) error {
	return dw.WriteSRAM(dw.MCU.SREG().Mem16(), []byte{b})
}

func (dw *DebugWIRE) GetSREG() (byte, error) {
	c := make([]byte, 1)
	if err := dw.ReadSRAM(dw.MCU.SREG().Mem16(), c); err != nil {
		return 0, err
	}
	return c[0], nil
}
