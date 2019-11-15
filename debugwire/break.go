package debugwire

import (
	"fmt"
)

func (dw *DebugWire) SendBreak() error {
	b, err := dw.Port.SendBreak()
	if err != nil {
		return err
	}

	if b != 0x55 {
		return fmt.Errorf("debugwire: bad break response. expected 0x55, got 0x%02x", b)
	}

	dw.afterBreak = true
	return nil
}

func (dw *DebugWire) RecvBreak() error {
	b, err := dw.Port.RecvBreak()
	if err != nil {
		return err
	}

	if b != 0x55 {
		return fmt.Errorf("debugwire: bad break received. expected 0x55, got 0x%02x", b)
	}

	dw.afterBreak = true
	return nil
}
