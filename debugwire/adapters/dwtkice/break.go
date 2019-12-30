package dwtkice

func (dw *DwtkIceAdapter) SendBreak() error {
	return dw.controlIn(cmdSendBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) RecvBreak() error {
	return dw.controlIn(cmdRecvBreak, 0, 0, nil)
}
