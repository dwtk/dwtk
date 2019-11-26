package debugwire

import (
	"fmt"
)

func (dw *DebugWire) SendBreak() error {
	if err := dw.device.SendBreak(); err != nil {
		return err
	}

	return dw.RecvBreak()
}

func (dw *DebugWire) RecvBreak() error {
	b, err := dw.device.RecvBreak()
	if err != nil {
		return err
	}

	if b != 0x55 {
		return fmt.Errorf("debugwire: bad break received. expected 0x55, got 0x%02x", b)
	}

	dw.afterBreak = true
	return nil
}
