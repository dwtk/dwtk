package debugwire

import (
	"errors"
	"strings"

	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWIRE) SetHwBreakpoint(addr uint16) bool {
	if dw.hwBreakpointSet {
		return false
	}

	dw.hwBreakpointSet = true
	dw.hwBreakpoint = addr
	return true
}

func (dw *DebugWIRE) ClearHwBreakpoint() {
	dw.hwBreakpointSet = false
	dw.hwBreakpoint = 0
}

func (dw *DebugWIRE) SetSwBreakpoint(addr uint16) error {
	f := make([]byte, 2)
	if err := dw.ReadFlash(addr, f); err != nil {
		return err
	}

	dw.swBreakpoints[addr] = (uint16(f[1]) << 8) | uint16(f[0])
	return dw.WriteFlashInstruction(addr, avr.BREAK())
}

func (dw *DebugWIRE) ClearSwBreakpoint(addr uint16) error {
	bp, ok := dw.swBreakpoints[addr]
	if !ok {
		return nil
	}

	if err := dw.WriteFlashInstruction(addr, bp); err != nil {
		return err
	}

	delete(dw.swBreakpoints, addr)
	return nil
}

func (dw *DebugWIRE) ClearSwBreakpoints() error {
	// this is used for recovery, so try to clear everything
	errs := []string{}
	for k := range dw.swBreakpoints {
		if err := dw.ClearSwBreakpoint(k); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func (dw *DebugWIRE) HasSwBreakpoints() bool {
	return len(dw.swBreakpoints) > 0
}
