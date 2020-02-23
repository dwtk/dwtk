package usbserial

import (
	"fmt"
)

func (us *UsbSerialAdapter) Disable() error {
	if err := us.device.Write([]byte{0x06}); err != nil {
		return err
	}

	fmt.Println("debugWIRE was disabled for target device, and it can be flashed using an SPI ISP now.")
	fmt.Println("this must be done without a target power cycle.")
	return nil
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

func (us *UsbSerialAdapter) ReadSignature() (uint16, error) {
	if err := us.device.Write([]byte{0xf3}); err != nil {
		return 0, err
	}

	return us.device.ReadWord()
}
