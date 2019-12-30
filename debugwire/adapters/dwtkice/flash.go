package dwtkice

func (dw *DwtkIceAdapter) WriteFlashPage(start uint16, data []byte) error {
	return dw.controlOut(cmdWriteFlashPage, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadFlash(start uint16, data []byte) error {
	return dw.controlIn(cmdReadFlash, start, 0, data)
}
