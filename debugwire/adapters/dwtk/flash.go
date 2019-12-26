package dwtk

func (dw *DwtkAdapter) ReadFlash(start uint16, data []byte) error {
	return dw.controlIn(cmdReadFlash, start, 0, data)
}
