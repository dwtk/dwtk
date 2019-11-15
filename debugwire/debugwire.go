package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/avr"
	"golang.rgm.io/dwtk/usbserial"
)

type DebugWire struct {
	Port *usbserial.UsbSerial
	MCU  *avr.MCU

	HwBreakpoint    uint16
	HwBreakpointSet bool
	Timers          bool

	afterBreak bool
}

func New(portDevice string, baudrate uint32) (*DebugWire, error) {
	u, err := usbserial.New(portDevice, baudrate)
	if err != nil {
		return nil, err
	}

	rv := &DebugWire{
		Port:            u,
		HwBreakpointSet: false,
		Timers:          true,
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
		return nil, fmt.Errorf("debugwire: failed to detect MCU from signature: %s", sign)
	}

	return rv, nil
}

func (dw *DebugWire) Close() error {
	return dw.Port.Close()
}
