package debugwire

import (
	"context"
	"fmt"

	"github.com/dwtk/dwtk/avr"
	"github.com/dwtk/dwtk/debugwire/adapters"
)

type DebugWIRE struct {
	MCU    *avr.MCU
	Timers bool

	adapter         adapters.Adapter
	hwBreakpoint    uint16
	hwBreakpointSet bool
	swBreakpoints   map[uint16]uint16
}

func New(device string, baudrate uint32) (*DebugWIRE, error) {
	a, err := adapters.New(device, baudrate)
	if err != nil {
		return nil, err
	}

	rv := &DebugWIRE{
		Timers: true,

		adapter:         a,
		hwBreakpointSet: false,
		swBreakpoints:   make(map[uint16]uint16, 1),
	}

	sign, err := a.GetSignature()
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

func (dw *DebugWIRE) Close() error {
	return dw.adapter.Close()
}

func (dw *DebugWIRE) Info() string {
	return dw.adapter.Info()
}

func (dw *DebugWIRE) Disable() error {
	return dw.adapter.Disable()
}

func (dw *DebugWIRE) Reset() error {
	return dw.adapter.Reset()
}

func (dw *DebugWIRE) GetSignature() (uint16, error) {
	return dw.adapter.GetSignature()
}

func (dw *DebugWIRE) SendBreak() error {
	return dw.adapter.SendBreak()
}

func (dw *DebugWIRE) RecvBreak() error {
	return dw.adapter.RecvBreak()
}

func (dw *DebugWIRE) Go() error {
	return dw.adapter.Go()
}

func (dw *DebugWIRE) Step() error {
	return dw.adapter.Step()
}

func (dw *DebugWIRE) Continue() error {
	return dw.adapter.Continue(dw.hwBreakpoint, dw.hwBreakpointSet, dw.Timers)
}

func (dw *DebugWIRE) Wait(ctx context.Context, c chan bool) error {
	return dw.adapter.Wait(ctx, c)
}

func (dw *DebugWIRE) WriteInstruction(inst uint16) error {
	return dw.adapter.WriteInstruction(inst)
}

func (dw *DebugWIRE) SetPC(pc uint16) error {
	return dw.adapter.SetPC(pc)
}

func (dw *DebugWIRE) GetPC() (uint16, error) {
	return dw.adapter.GetPC()
}

func (dw *DebugWIRE) WriteRegisters(start byte, regs []byte) error {
	return dw.adapter.WriteRegisters(start, regs)
}

func (dw *DebugWIRE) ReadRegisters(start byte, regs []byte) error {
	return dw.adapter.ReadRegisters(start, regs)
}

func (dw *DebugWIRE) WriteSRAM(start uint16, data []byte) error {
	return dw.adapter.WriteSRAM(start, data)
}

func (dw *DebugWIRE) ReadSRAM(start uint16, data []byte) error {
	return dw.adapter.ReadSRAM(start, data)
}

func (dw *DebugWIRE) ReadFuses() ([]byte, error) {
	return dw.adapter.ReadFuses()
}
