package usbserial

import (
	"fmt"

	"golang.rgm.io/dwtk/usbserial"
)

type UsbSerialAdapter struct {
	device     *usbserial.UsbSerial
	serialPort string
	baudrate   uint32
	afterBreak bool
}

func New(serialPort string, baudrate uint32) (*UsbSerialAdapter, error) {
	var err error
	if serialPort == "" {
		serialPort, err = guessSerialPort()
		if err != nil {
			return nil, err
		}
	}

	if baudrate == 0 {
		var err error
		baudrate, err = guessBaudrate(serialPort)
		if err != nil {
			return nil, err
		}
	}

	u, err := usbserial.Open(serialPort, baudrate)
	if err != nil {
		return nil, err
	}

	return &UsbSerialAdapter{
		device:     u,
		serialPort: serialPort,
		baudrate:   baudrate,
		afterBreak: false,
	}, nil
}

func (us *UsbSerialAdapter) Close() error {
	return us.device.Close()
}

func (us *UsbSerialAdapter) Info() string {
	return fmt.Sprintf("Serial Port (USB Serial): %s\nBaud Rate: %d bps\n", us.serialPort, us.baudrate)
}
