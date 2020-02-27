package usbserial

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dwtk/dwtk/avr"
	"github.com/dwtk/dwtk/internal/logger"
	"github.com/dwtk/dwtk/internal/usbserial"
)

var (
	errNotSupported = errors.New("debugwire: usbserial: operation not supported by Usb Serial")
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
		logger.Debug.Printf(" * Detected baudrate: %d", baudrate)
	}

	u, err := usbserial.Open(serialPort, baudrate)
	if err != nil {
		return nil, err
	}

	rv := &UsbSerialAdapter{
		device:     u,
		serialPort: serialPort,
		baudrate:   baudrate,
		afterBreak: false,
	}

	if err := rv.SendBreak(); err != nil {
		rv.Close()
		return nil, err
	}

	return rv, nil
}

func (us *UsbSerialAdapter) Close() error {
	return us.device.Close()
}

func (us *UsbSerialAdapter) Info() string {
	return fmt.Sprintf("Serial Port (USB Serial): %s\nBaud Rate: %d bps\n", us.serialPort, us.baudrate)
}

func (us *UsbSerialAdapter) SetMCU(mcu *avr.MCU) {
	// we don't need the mcu for anything
}

func (us *UsbSerialAdapter) Enable() error {
	return errNotSupported
}

func (us *UsbSerialAdapter) ChipErase() error {
	return errNotSupported
}

func (us *UsbSerialAdapter) WriteLFuse(data byte) error {
	return errNotSupported
}

func (us *UsbSerialAdapter) WriteHFuse(data byte) error {
	return errNotSupported
}

func (us *UsbSerialAdapter) WriteEFuse(data byte) error {
	return errNotSupported
}

func (us *UsbSerialAdapter) WriteLock(data byte) error {
	return errNotSupported
}
