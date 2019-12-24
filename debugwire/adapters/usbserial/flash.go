package usbserial

func (us *UsbSerialAdapter) ReadFlash(start uint16, b []byte) error {
	c := []byte{
		byte(start), byte(start >> 8), // Z
	}
	if err := us.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16(len(b) * 2)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x00,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x02,
		0x20,
	}
	if err := us.device.Write(c); err != nil {
		return err
	}

	return us.device.Read(b)
}
