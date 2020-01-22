package dwtkice

import (
	"fmt"

	"github.com/dwtk/dwtk/internal/logger"
	"github.com/dwtk/dwtk/internal/usbfs"
)

const (
	VID = 0x1d50 // OpenMoko, Inc.
	PID = 0x614c // dwtk In-Circuit Emulator
)

const (
	cmdGetError = iota + 0x40
	cmdGetBaudratePrescaler
	cmdDetectBaudrate
	cmdSetBaudrate
	cmdDisable
	cmdReset
	cmdGetSignature
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

var (
	cmds = map[byte]string{
		cmdGetError:             "cmdGetError",
		cmdGetBaudratePrescaler: "cmdGetBaudratePrescaler",
		cmdDetectBaudrate:       "cmdDetectBaudrate",
		cmdSetBaudrate:          "cmdSetBaudrate",
		cmdDisable:              "cmdDisable",
		cmdReset:                "cmdReset",
		cmdGetSignature:         "cmdGetSignature",
		cmdSendBreak:            "cmdSendBreak",
		cmdRecvBreak:            "cmdRecvBreak",
		cmdGo:                   "cmdGo",
		cmdStep:                 "cmdStep",
		cmdContinue:             "cmdContinue",
		cmdWait:                 "cmdWait",
		cmdWriteInstruction:     "cmdWriteInstruction",
		cmdSetPC:                "cmdSetPC",
		cmdGetPC:                "cmdGetPC",
		cmdRegisters:            "cmdRegisters",
		cmdSRAM:                 "cmdSRAM",
		cmdReadFlash:            "cmdReadFlash",
		cmdWriteFlashPage:       "cmdWriteFlashPage",
		cmdReadFuses:            "cmdReadFuses",
	}
)

type DwtkIceAdapter struct {
	device     *usbfs.Device
	ubrr       uint16
	baudrate   uint32
	afterBreak bool
}

func New(baudrate uint32) (*DwtkIceAdapter, error) {
	devices, err := usbfs.GetDevices(VID, PID)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}
	if len(devices) > 1 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: more than one dwtk-ice device found. this is not supported")
	}

	rv := &DwtkIceAdapter{
		device:     devices[0],
		afterBreak: false,
	}
	if err := rv.device.Open(); err != nil {
		return nil, err
	}
	logger.Debug.Printf(" * Detected dwtk-ice %s", rv.device.GetVersion())

	if baudrate == 0 {
		rv.ubrr, err = rv.detectBaudrate()
		if err != nil {
			rv.Close()
			return nil, err
		}
	} else {
		rv.ubrr, err = rv.baudrateToUbrr(baudrate)
		if err != nil {
			rv.Close()
			return nil, err
		}

		if err := rv.setBaudrate(rv.ubrr); err != nil {
			rv.Close()
			return nil, err
		}
	}

	rv.baudrate, err = rv.ubrrToBaudrate(rv.ubrr)
	if err != nil {
		rv.Close()
		return nil, err
	}
	logger.Debug.Printf(" * Actual baudrate: %d", rv.baudrate)

	return rv, nil
}

func (dw *DwtkIceAdapter) Close() error {
	return dw.device.Close()
}

func (dw *DwtkIceAdapter) Info() string {
	return fmt.Sprintf("dwtk-ice %s\n\nBaud Rate: %d bps\nBaud Rate Register: 0x%04x\n",
		dw.device.GetVersion(),
		dw.baudrate,
		dw.ubrr,
	)
}
