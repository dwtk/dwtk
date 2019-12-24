package usbserial

func (us *UsbSerialAdapter) Disable() error {
	return us.device.Write([]byte{0x06})
}

func (us *UsbSerialAdapter) Reset() error {
	if err := us.SendBreak(); err != nil {
		return err
	}

	if err := us.device.Write([]byte{0x07}); err != nil {
		return err
	}

	return us.RecvBreak()
}

func (us *UsbSerialAdapter) GetSignature() (uint16, error) {
	if err := us.device.Write([]byte{0xf3}); err != nil {
		return 0, err
	}

	return us.device.ReadWord()
}
