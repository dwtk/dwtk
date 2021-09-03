package dwtkice

import (
	"fmt"
	"strings"
	"time"

	"github.com/dwtk/dwtk/internal/logger"
	"github.com/rafaelmartins/usbfs"
)

const (
	vid = 0x1d50 // OpenMoko, Inc.
	pid = 0x614c // dwtk In-Circuit Emulator
)

const (
	capDw = (1 << iota)
	capSpi
)

const (
	cmdGetError = iota + 1
	cmdGetCapabilities
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
	errNone = iota
	errUnsupported
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
		cmdGetError:        "cmdGetError",
		cmdGetCapabilities: "cmdGetCapabilities",

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

	iceErrors = map[byte]func(byte, byte) error{
		errUnsupported: func(_ byte, _ byte) error {
			return errCmdUnsupported
		},
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

	if e[0] == errNone {
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
	spi bool
}

func newDevice(serialNumber string) (*device, error) {
	serials := []string{}
	devices, err := usbfs.List(func(d *usbfs.Device) bool {
		idVendor, err := d.IdVendor()
		if err != nil || idVendor != vid {
			return false
		}

		idProduct, err := d.IdProduct()
		if err != nil || idProduct != pid {
			return false
		}

		if serialNumber != "" {
			serial, err := d.Serial()
			if err != nil || serial != serialNumber {
				return false
			}
			serials = append(serials, serial)
		}

		return true
	})
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		if serialNumber != "" {
			return nil, fmt.Errorf("debugwire: dwtk-ice: device not found: %s", serialNumber)
		}
		return nil, nil
	}
	if len(devices) > 1 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: more than one device found. this is not supported: %s",
			strings.Join(serials, ", "))
	}
	dev := devices[0]

	if err := dev.Open(); err != nil {
		return nil, err
	}
	rv := &device{
		dev: dev,
		spi: false,
	}
	b := make([]byte, 1)
	if err := rv.controlIn(cmdGetCapabilities, 0, 0, b); err != nil {
		return nil, err
	}
	if b[0]&capDw == 0 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: debugwire not supported. this is probably a connection problem")
	}
	if b[0]&capSpi != 0 {
		rv.spi = true
	}
	return rv, nil
}

func (d *device) close() error {
	if d.dev != nil {
		return d.dev.Close()
	}
	return nil
}

func (d *device) getVersion() string {
	bcdDevice, err := d.dev.BcdDevice()
	if err != nil {
		return "UNKNOWN"
	}
	return fmt.Sprintf("%x.%02x", bcdDevice/0x0100, bcdDevice%0x0100)
}

func (d *device) getSerial() string {
	rv, _ := d.dev.Serial()
	return rv
}

func (d *device) controlGetError() error {
	f := make([]byte, 3)
	if err := d.dev.Control(usbfs.RequestTypeVendor, usbfs.RequestRecipientDevice, usbfs.DirectionIn, cmdGetError, 0, 0, f, 5*time.Second); err != nil {
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
	if err := d.dev.Control(usbfs.RequestTypeVendor, usbfs.RequestRecipientDevice, usbfs.DirectionIn, req, val, idx, f, 5*time.Second); err != nil {
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
	if err := d.dev.Control(usbfs.RequestTypeVendor, usbfs.RequestRecipientDevice, usbfs.DirectionOut, req, val, idx, data, 5*time.Second); err != nil {
		return err
	}
	for _, c := range data {
		logger.Debug.Printf(">>> 0x%02x", c)
	}
	return d.controlGetError()
}
