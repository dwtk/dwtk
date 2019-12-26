package dwtk

import (
	"errors"
	"fmt"

	"golang.rgm.io/dwtk/logger"
)

var (
	dwtkErrors = map[uint8]error{
		1: errors.New("debugwire: dwtk: hardware not initialized"),
		2: errors.New("debugwire: dwtk: baudrate detection failed"),
		3: errors.New("debugwire: dwtk: got unexpected byte echoed back"),
		4: errors.New("debugwire: dwtk: got unexpected break value"),
	}
)

func (dw *DwtkAdapter) controlGetError() error {
	f := make([]byte, 1)
	if err := dw.device.ControlIn(cmdGetError, 0, 0, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< cmdGetError: 0x%02x", f[0])
	if f[0] == 0 {
		return nil
	}
	err, ok := dwtkErrors[f[0]]
	if !ok {
		return fmt.Errorf("debugwire: dwtk: unrecognized hardware error: 0x%02x", f[0])
	}
	return err
}

func (dw *DwtkAdapter) controlIn(req byte, val uint16, idx uint16, data []byte) error {
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

func (dw *DwtkAdapter) controlOut(req byte, val uint16, idx uint16, data []byte) error {
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
