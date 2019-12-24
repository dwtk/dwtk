package usbserial

func (us *UsbSerialAdapter) WriteInstruction(inst uint16) error {
	c := []byte{
		0x64,
		0xd2, byte(inst >> 8), byte(inst),
		0x23,
	}
	return us.device.Write(c)
}
