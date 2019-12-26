package dwtk

func (dw *DwtkAdapter) WriteSRAM(start uint16, data []byte) error {
	return dw.controlOut(cmdSRAM, start, 0, data)
}

func (dw *DwtkAdapter) ReadSRAM(start uint16, data []byte) error {
	return dw.controlIn(cmdSRAM, start, 0, data)
}
