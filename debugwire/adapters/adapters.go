package adapters

import (
	"context"
	"fmt"

	"github.com/dwtk/dwtk/avr"
	"github.com/dwtk/dwtk/debugwire/adapters/dwtkice"
	"github.com/dwtk/dwtk/debugwire/adapters/usbserial"
)

type Adapter interface {
	Close() error
	Info() string
	SetMCU(mcu *avr.MCU)

	Enable() error
	Disable() error
	Reset() error
	ReadSignature() (uint16, error)
	ChipErase() error

	SendBreak() error
	RecvBreak() error

	Go() error
	ResetAndGo() error
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
	WriteFlashPage(start uint16, data []byte) error
	EraseFlashPage(start uint16) error

	ReadFuses() ([]byte, error)
	WriteLFuse(data byte) error
	WriteHFuse(data byte) error
	WriteEFuse(data byte) error
	WriteLock(data byte) error
}

func New(dwtkIce bool, serialPort string, baudrate uint32) (Adapter, error) {
	if dwtkIce || serialPort == "" {
		adapter, err := dwtkice.New()
		if err != nil {
			return nil, err
		}
		if adapter != nil {
			return adapter, nil
		}
		if dwtkIce {
			return nil, fmt.Errorf("debugwire: adapters: dwtk-ice requested but no device found")
		}
	}

	adapter, err := usbserial.New(serialPort, baudrate)
	if err != nil {
		return nil, err
	}
	if adapter == nil {
		return nil, fmt.Errorf("debugwire: adapters: no device found")
	}
	return adapter, nil
}
