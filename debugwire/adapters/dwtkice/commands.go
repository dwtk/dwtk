package dwtkice

func (dw *DwtkIceAdapter) Disable() error {
	return dw.controlIn(cmdDisable, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Reset() error {
	return dw.controlIn(cmdReset, 0, 0, nil)
}
