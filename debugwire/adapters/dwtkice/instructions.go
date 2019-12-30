package dwtkice

func (dw *DwtkIceAdapter) WriteInstruction(inst uint16) error {
	return dw.controlIn(cmdWriteInstruction, inst, 0, nil)
}
