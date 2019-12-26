package dwtk

func (dw *DwtkAdapter) WriteInstruction(inst uint16) error {
	return dw.controlIn(cmdWriteInstruction, inst, 0, nil)
}
