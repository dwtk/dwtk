package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/avr"
	"golang.rgm.io/dwtk/usbserial"
)

type DebugWire struct {
	MCU    *avr.MCU
	Timers bool

	device          *usbserial.UsbSerial
	hwBreakpoint    uint16
	hwBreakpointSet bool
	swBreakpoints   map[uint16]uint16
	afterBreak      bool
}

func New(device string, baudrate uint32) (*DebugWire, error) {
	u, err := usbserial.Open(device, baudrate)
	if err != nil {
		return nil, err
	}

	rv := &DebugWire{
		Timers: true,

		device:          u,
		hwBreakpointSet: false,
		swBreakpoints:   make(map[uint16]uint16, 1),
		afterBreak:      false,
	}

	sign, err := rv.GetSignature()
	if err != nil {
		rv.Close()
		return nil, err
	}

	rv.MCU = avr.GetMCU(sign)
	if rv.MCU == nil {
		rv.Close()
		return nil, fmt.Errorf(`debugwire: failed to detect MCU from signature: 0x%04x
Please open an issue/pull request: https://github.com/rafaelmartins/dwtk`, sign)
	}

	return rv, nil
}

func (dw *DebugWire) Close() error {
	return dw.device.Close()
}
