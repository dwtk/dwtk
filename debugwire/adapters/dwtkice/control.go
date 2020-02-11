package dwtkice

import (
	"fmt"

	"github.com/dwtk/dwtk/internal/logger"
)

func (dw *DwtkIceAdapter) codeToError(e []byte) error {
	if len(e) < 3 {
		return fmt.Errorf("debugwire: dwtk-ice: invalid error: %v", e)
	}

	if e[0] == 0 {
		return nil
	}

	errFunc, ok := iceErrors[e[0]]
	if !ok {
		return fmt.Errorf("debugwire: dwtk-ice: unrecognized hardware error: 0x%02x", e[0])
	}
	return errFunc(e[1], e[2])
}

func (dw *DwtkIceAdapter) controlGetError() error {
	f := make([]byte, 3)
	if err := dw.device.ControlIn(cmdGetError, 0, 0, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< cmdGetError: 0x%02x -> [0x%02x, 0x%02x]", f[0], f[1], f[2])
	return dw.codeToError(f)
}

func (dw *DwtkIceAdapter) controlIn(req byte, val uint16, idx uint16, data []byte) error {
	cmd, ok := cmds[req]
	if ok {
		logger.Debug.Printf("<<< %s(0x%04x, 0x%04x)", cmd, val, idx)
	} else {
		logger.Debug.Printf("<<< %d(0x%04x, 0x%04x)", req, val, idx)
	}
	f := make([]byte, len(data)+3)
	if err := dw.device.ControlIn(req, val, idx, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< error: 0x%02x -> [0x%02x, 0x%02x]", f[0], f[1], f[2])
	for i, d := range f[3:] {
		data[i] = d
		logger.Debug.Printf("<<< 0x%02x", d)
	}
	return dw.codeToError(f)
}

func (dw *DwtkIceAdapter) controlOut(req byte, val uint16, idx uint16, data []byte) error {
	cmd, ok := cmds[req]
	if ok {
		logger.Debug.Printf(">>> %s(0x%04x, 0x%04x)", cmd, val, idx)
	} else {
		logger.Debug.Printf(">>> %d(0x%04x, 0x%04x)", req, val, idx)
	}
	if err := dw.device.ControlOut(req, val, idx, data); err != nil {
		return err
	}
	for _, d := range data {
		logger.Debug.Printf(">>> 0x%02x", d)
	}
	return dw.controlGetError()
}
