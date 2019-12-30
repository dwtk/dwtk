package dwtkice

func (dw *DwtkIceAdapter) WriteSRAM(start uint16, data []byte) error {
	return dw.controlOut(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadSRAM(start uint16, data []byte) error {
	return dw.controlIn(cmdSRAM, start, 0, data)
}
