package debugwire

func (dw *DebugWIRE) Go() error {
	c := []byte{
		0x40,
		0x30,
	}
	if err := dw.device.Write(c); err != nil {
		return err
	}
	return dw.device.Commit()
}

func (dw *DebugWIRE) Step() error {
	c := []byte{
		0x60,
		0x31,
	}
	if err := dw.device.Write(c); err != nil {
		return err
	}

	return dw.RecvBreak()
}

func (dw *DebugWIRE) Continue() error {
	c := []byte{}
	t := byte(0x60)
	if dw.hwBreakpointSet {
		bp := dw.hwBreakpoint / 2
		c = append(c, 0xd1, byte(bp>>8), byte(bp))
		t = 0x61
	}
	if dw.Timers {
		t -= 0x20
	}
	if err := dw.device.Write(append(c, t, 0x30)); err != nil {
		return err
	}
	return dw.device.Commit()
}
