package usbserial

import (
	"fmt"
	"strings"

	"golang.rgm.io/dwtk/internal/usbserial"
	"golang.rgm.io/dwtk/logger"
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
		devices, err := usbserial.ListDevices()
		if err != nil {
			return nil, err
		}

		if len(devices) == 0 {
			return nil, nil
		}
		if len(devices) > 1 {
			return nil, fmt.Errorf("debugwire: usbserial: more than one dwtk device found. this is not supported: %s",
				strings.Join(devices, ", "))
		}

		logger.Debug.Printf(" * Detected serial port: %s", devices[0])
		serialPort = devices[0]
	}

	if baudrate == 0 {
		var err error
		baudrate, err = detectBaudrate(serialPort)
		if err != nil {
			return nil, err
		}
		logger.Debug.Printf(" * Detected baudrate: %d\n", baudrate)
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
