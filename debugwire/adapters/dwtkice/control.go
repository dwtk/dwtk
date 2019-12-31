package dwtkice

import (
	"errors"
	"fmt"

	"github.com/dwtk/dwtk/internal/logger"
)

var (
	iceErrors = map[uint8]error{
		1: errors.New("debugwire: dwtk-ice: hardware not initialized"),
		2: errors.New("debugwire: dwtk-ice: baudrate detection failed"),
		3: errors.New("debugwire: dwtk-ice: got unexpected byte echoed back"),
		4: errors.New("debugwire: dwtk-ice: got unexpected break value"),
	}
)

func (dw *DwtkIceAdapter) controlGetError() error {
	f := make([]byte, 1)
	if err := dw.device.ControlIn(cmdGetError, 0, 0, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< cmdGetError: 0x%02x", f[0])
	if f[0] == 0 {
		return nil
	}
	err, ok := iceErrors[f[0]]
	if !ok {
		return fmt.Errorf("debugwire: dwtk-ice: unrecognized hardware error: 0x%02x", f[0])
	}
	return err
}

func (dw *DwtkIceAdapter) controlIn(req byte, val uint16, idx uint16, data []byte) error {
	cmd, ok := cmds[req]
	if ok {
		logger.Debug.Printf("<<< %s(0x%04x, 0x%04x)", cmd, val, idx)
	} else {
		logger.Debug.Printf("<<< %d(0x%04x, 0x%04x)", req, val, idx)
	}
	if err := dw.device.ControlIn(req, val, idx, data); err != nil {
		return err
	}
	for _, d := range data {
		logger.Debug.Printf("<<< 0x%02x", d)
	}
	return dw.controlGetError()
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
