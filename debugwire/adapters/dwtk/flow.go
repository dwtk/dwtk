package dwtk

func (dw *DwtkAdapter) Go() error {
	return dw.controlIn(cmdGo, 0, 0, nil)
}

func (dw *DwtkAdapter) Step() error {
	return dw.controlIn(cmdStep, 0, 0, nil)
}

func (dw *DwtkAdapter) Continue(hwBreakpoint uint16, hwBreakpointSet bool, timers bool) error {
	// idx: byte 0 -> hw bp set
	//      byte 1 -> timers
	idx := uint16(0)
	if hwBreakpointSet {
		idx |= (1 << 0)
	}
	if timers {
		idx |= (1 << 1)
	}
	return dw.controlIn(cmdContinue, hwBreakpoint, idx, nil)
}
