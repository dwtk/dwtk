package dwtkice

func (dw *DwtkIceAdapter) WriteFlashPage(start uint16, data []byte) error {
	return dw.controlOut(cmdWriteFlashPage, start, 0, data)
}

func (dw *DwtkIceAdapter) EraseFlashPage(start uint16) error {
	return dw.controlIn(cmdEraseFlashPage, start, 0, nil)
}

func (dw *DwtkIceAdapter) ReadFlash(start uint16, data []byte) error {
	return dw.controlIn(cmdReadFlash, start, 0, data)
}
