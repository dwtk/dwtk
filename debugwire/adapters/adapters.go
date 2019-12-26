package adapters

import (
	"context"

	"golang.rgm.io/dwtk/debugwire/adapters/dwtk"
	"golang.rgm.io/dwtk/debugwire/adapters/usbserial"
)

type Adapter interface {
	Close() error

	Info() string

	Disable() error
	Reset() error
	GetSignature() (uint16, error)

	SendBreak() error
	RecvBreak() error

	Go() error
	Step() error
	Continue(hwBreakpoint uint16, hwBreakpointSet bool, timers bool) error

	Wait(ctx context.Context, c chan bool) error

	WriteInstruction(inst uint16) error

	SetPC(pc uint16) error
	GetPC() (uint16, error)
	WriteRegisters(start byte, regs []byte) error
	ReadRegisters(start byte, regs []byte) error

	WriteSRAM(start uint16, data []byte) error
	ReadSRAM(start uint16, data []byte) error

	ReadFlash(start uint16, data []byte) error
}

func New(serialPort string, baudrate uint32) (Adapter, error) {
	if serialPort == "" {
		adapter, err := dwtk.New(baudrate)
		if err != nil {
			return nil, err
		}

		if adapter != nil {
			return adapter, nil
		}
	}
	return usbserial.New(serialPort, baudrate)
}
