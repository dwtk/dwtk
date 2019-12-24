package usbserial

func (us *UsbSerialAdapter) WriteSRAM(start uint16, data []byte) error {
	c := []byte{
		byte(start), byte(start >> 8),
	}
	if err := us.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16((len(data) * 2) + 1)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x01,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x04,
		0x20,
	}
	return us.device.Write(append(c, data...))
}

func (us *UsbSerialAdapter) ReadSRAM(start uint16, data []byte) error {
	c := []byte{
		byte(start), byte(start >> 8),
	}
	if err := us.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16((len(data) * 2) + 1)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x00,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x00,
		0x20,
	}
	if err := us.device.Write(c); err != nil {
		return err
	}

	return us.device.Read(data)
}
