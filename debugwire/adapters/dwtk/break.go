package dwtk

func (dw *DwtkAdapter) SendBreak() error {
	return dw.controlIn(cmdSendBreak, 0, 0, nil)
}

func (dw *DwtkAdapter) RecvBreak() error {
	return dw.controlIn(cmdRecvBreak, 0, 0, nil)
}
