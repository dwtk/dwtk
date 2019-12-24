package usbserial

func (us *UsbSerialAdapter) Go() error {
	c := []byte{
		0x40,
		0x30,
	}
	if err := us.device.Write(c); err != nil {
		return err
	}
	return us.device.Commit()
}

func (us *UsbSerialAdapter) Step() error {
	c := []byte{
		0x60,
		0x31,
	}
	if err := us.device.Write(c); err != nil {
		return err
	}

	return us.RecvBreak()
}

func (us *UsbSerialAdapter) Continue(hwBreakpoint uint16, hwBreakpointSet bool, timers bool) error {
	c := []byte{}
	t := byte(0x60)
	if hwBreakpointSet {
		bp := hwBreakpoint / 2
		c = append(c, 0xd1, byte(bp>>8), byte(bp))
		t = 0x61
	}
	if timers {
		t -= 0x20
	}
	if err := us.device.Write(append(c, t, 0x30)); err != nil {
		return err
	}
	return us.device.Commit()
}
