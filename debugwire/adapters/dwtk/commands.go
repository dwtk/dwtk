package dwtk

func (dw *DwtkAdapter) Disable() error {
	return dw.controlIn(cmdDisable, 0, 0, nil)
}

func (dw *DwtkAdapter) Reset() error {
	return dw.controlIn(cmdReset, 0, 0, nil)
}

func (dw *DwtkAdapter) GetSignature() (uint16, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetSignature, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}
