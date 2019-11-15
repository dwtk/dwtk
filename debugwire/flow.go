package debugwire

func (dw *DebugWire) Go() error {
	c := []byte{
		0x40,
		0x30,
	}
	if err := dw.Port.Write(c); err != nil {
		return err
	}
	return dw.Port.Commit()
}

func (dw *DebugWire) Step() error {
	c := []byte{
		0x60,
		0x31,
	}
	if err := dw.Port.Write(c); err != nil {
		return err
	}

	return dw.RecvBreak()
}

func (dw *DebugWire) Continue() error {
	c := []byte{}
	t := byte(0x60)
	if dw.HwBreakpointSet {
		c = append(c, 0xd1, byte(dw.HwBreakpoint>>8), byte(dw.HwBreakpoint))
		t = 0x61
	}
	if dw.Timers {
		t -= 0x20
	}
	if err := dw.Port.Write(append(c, t, 0x30)); err != nil {
		return err
	}
	return dw.Port.Commit()
}
