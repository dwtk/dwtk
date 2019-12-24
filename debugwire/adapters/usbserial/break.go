package usbserial

import (
	"fmt"
)

func (us *UsbSerialAdapter) SendBreak() error {
	if err := us.device.SendBreak(); err != nil {
		return err
	}

	return us.RecvBreak()
}

func (us *UsbSerialAdapter) RecvBreak() error {
	b, err := us.device.RecvBreak()
	if err != nil {
		return err
	}

	if b != 0x55 {
		return fmt.Errorf("debugwire: bad break received. expected 0x55, got 0x%02x", b)
	}

	us.afterBreak = true
	return nil
}
