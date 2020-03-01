package dwtkice

import (
	"fmt"

	"github.com/dwtk/dwtk/internal/logger"
	"github.com/dwtk/dwtk/internal/usbfs"
)

const (
	vid = 0x1d50 // OpenMoko, Inc.
	pid = 0x614c // dwtk In-Circuit Emulator
)

const (
	cmdGetError = iota + 1
)

const (
	cmdSpiPgmEnable = iota + 0x20
	cmdSpiPgmDisable
	cmdSpiCommand
	cmdSpiReset
)

const (
	cmdDetectBaudrate = iota + 0x40
	cmdGetBaudrate
	cmdDisable
	cmdReset
	cmdReadSignature
	cmdSendBreak
	cmdRecvBreak
	cmdGo
	cmdStep
	cmdContinue
	cmdWait
	cmdWriteInstruction
	cmdSetPC
	cmdGetPC
	cmdRegisters
	cmdSRAM
	cmdReadFlash
	cmdWriteFlashPage
	cmdEraseFlashPage
	cmdReadFuses
)

const (
	errSpiPgmEnable = iota + 0x20
	errSpiEchoMismatch
)

const (
	errBaudrateDetection = iota + 0x40
	errEchoMismatch
	errBreakMismatch
	errTooLarge
)

var (
	cmds = map[byte]string{
		cmdGetError: "cmdGetError",

		cmdSpiPgmEnable:  "cmdSpiPgmEnable",
		cmdSpiPgmDisable: "cmdSpiPgmDisable",
		cmdSpiCommand:    "cmdSpiCommand",
		cmdSpiReset:      "cmdSpiReset",

		cmdDetectBaudrate:   "cmdDetectBaudrate",
		cmdGetBaudrate:      "cmdGetBaudrate",
		cmdDisable:          "cmdDisable",
		cmdReset:            "cmdReset",
		cmdReadSignature:    "cmdReadSignature",
		cmdSendBreak:        "cmdSendBreak",
		cmdRecvBreak:        "cmdRecvBreak",
		cmdGo:               "cmdGo",
		cmdStep:             "cmdStep",
		cmdContinue:         "cmdContinue",
		cmdWait:             "cmdWait",
		cmdWriteInstruction: "cmdWriteInstruction",
		cmdSetPC:            "cmdSetPC",
		cmdGetPC:            "cmdGetPC",
		cmdRegisters:        "cmdRegisters",
		cmdSRAM:             "cmdSRAM",
		cmdReadFlash:        "cmdReadFlash",
		cmdWriteFlashPage:   "cmdWriteFlashPage",
		cmdEraseFlashPage:   "cmdEraseFlashPage",
		cmdReadFuses:        "cmdReadFuses",
	}

	iceErrors = map[uint8]func(byte, byte) error{
		errSpiPgmEnable: func(_ byte, _ byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: SPI programming enable failed")
		},
		errSpiEchoMismatch: func(exp byte, got byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: got unexpected byte echoed back via SPI: expected 0x%02x, got 0x%02x", exp, got)
		},
		errBaudrateDetection: func(_ byte, _ byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: baudrate detection failed")
		},
		errEchoMismatch: func(exp byte, got byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: got unexpected byte echoed back: expected 0x%02x, got 0x%02x", exp, got)
		},
		errBreakMismatch: func(got byte, _ byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: got unexpected break value: expected 0x55, got 0x%02x", got)
		},
		errTooLarge: func(_ byte, _ byte) error {
			return fmt.Errorf("debugwire: dwtk-ice: read/write data is too large")
		},
	}
)

func codeToError(e []byte) error {
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

type device struct {
	dev *usbfs.Device
}

func newDevice() (*device, error) {
	devices, err := usbfs.GetDevices(vid, pid)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}
	if len(devices) > 1 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: more than one dwtk-ice device found. this is not supported")
	}
	if err := devices[0].Open(); err != nil {
		return nil, err
	}
	return &device{dev: devices[0]}, nil
}

func (d *device) close() error {
	return d.dev.Close()
}

func (d *device) getVersion() string {
	return fmt.Sprintf("%s", d.dev.GetVersion())
}

func (d *device) controlGetError() error {
	f := make([]byte, 3)
	if err := d.dev.ControlIn(cmdGetError, 0, 0, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< cmdGetError: 0x%02x -> [0x%02x, 0x%02x]", f[0], f[1], f[2])
	return codeToError(f)
}

func (d *device) controlIn(req byte, val uint16, idx uint16, data []byte) error {
	cmd, ok := cmds[req]
	if ok {
		logger.Debug.Printf("<<< %s(0x%04x, 0x%04x)", cmd, val, idx)
	} else {
		logger.Debug.Printf("<<< %d(0x%04x, 0x%04x)", req, val, idx)
	}
	f := make([]byte, len(data)+3)
	if err := d.dev.ControlIn(req, val, idx, f); err != nil {
		return err
	}
	logger.Debug.Printf("<<< error: 0x%02x -> [0x%02x, 0x%02x]", f[0], f[1], f[2])
	for i, c := range f[3:] {
		data[i] = c
		logger.Debug.Printf("<<< 0x%02x", c)
	}
	return codeToError(f)
}

func (d *device) controlOut(req byte, val uint16, idx uint16, data []byte) error {
	cmd, ok := cmds[req]
	if ok {
		logger.Debug.Printf(">>> %s(0x%04x, 0x%04x)", cmd, val, idx)
	} else {
		logger.Debug.Printf(">>> %d(0x%04x, 0x%04x)", req, val, idx)
	}
	if err := d.dev.ControlOut(req, val, idx, data); err != nil {
		return err
	}
	for _, c := range data {
		logger.Debug.Printf(">>> 0x%02x", c)
	}
	return d.controlGetError()
}
