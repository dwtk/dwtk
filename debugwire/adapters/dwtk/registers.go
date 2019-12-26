package dwtk

func (dw *DwtkAdapter) WriteRegisters(start byte, regs []byte) error {
	return dw.controlOut(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkAdapter) ReadRegisters(start byte, regs []byte) error {
	return dw.controlIn(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkAdapter) SetPC(pc uint16) error {
	return dw.controlIn(cmdSetPC, pc, 0, nil)
}

func (dw *DwtkAdapter) GetPC() (uint16, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetPC, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}
