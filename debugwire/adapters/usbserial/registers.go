package usbserial

func (us *UsbSerialAdapter) WriteRegisters(start byte, regs []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(regs)),
		0xc2, 0x05,
		0x20,
	}
	return us.device.Write(append(c, regs...))
}

func (us *UsbSerialAdapter) ReadRegisters(start byte, regs []byte) error {
	c := []byte{
		0x66,
		0xd0, 0x00, start, // ignoring high byte because registers are 0-31
		0xd1, 0x00, start + byte(len(regs)),
		0xc2, 0x01,
		0x20,
	}
	if err := us.device.Write(c); err != nil {
		return err
	}
	return us.device.Read(regs)
}

func (us *UsbSerialAdapter) SetPC(b uint16) error {
	us.afterBreak = false
	b /= 2
	c := []byte{
		0xd0, byte(b >> 8), byte(b),
	}
	return us.device.Write(c)
}

func (us *UsbSerialAdapter) GetPC() (uint16, error) {
	if err := us.device.Write([]byte{0xf0}); err != nil {
		return 0, err
	}

	rv, err := us.device.ReadWord()
	if err != nil {
		return 0, err
	}

	if us.afterBreak {
		if rv > 0 {
			rv--
		}
		us.afterBreak = false
	}

	rv *= 2
	return rv, nil
}
