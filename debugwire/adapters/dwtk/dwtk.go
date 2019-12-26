package dwtk

import (
	"fmt"

	"golang.rgm.io/dwtk/internal/usbfs"
)

const (
	VID          = 0x16c0
	PID          = 0x05dc
	Manufacturer = "dwtk.rgm.io"
	Product      = "dwtk-hardware"
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
	}
)

type DwtkAdapter struct {
	device     *usbfs.Device
	ubrr       uint16
	baudrate   uint32
	afterBreak bool
}

func New(baudrate uint32) (*DwtkAdapter, error) {
	devices, err := usbfs.GetDevices(VID, PID, Manufacturer, Product)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}
	if len(devices) > 1 {
		return nil, fmt.Errorf("debugwire: dwtk: more than one dwtk device found. this is not supported")
	}

	rv := &DwtkAdapter{
		device:     devices[0],
		afterBreak: false,
	}
	if err := rv.device.Open(); err != nil {
		return nil, err
	}

	if baudrate == 0 {
		rv.ubrr, err = rv.DetectBaudrate()
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

		if err := rv.SetBaudrate(rv.ubrr); err != nil {
			rv.Close()
			return nil, err
		}
	}

	rv.baudrate, err = rv.ubrrToBaudrate(rv.ubrr)
	if err != nil {
		rv.Close()
		return nil, err
	}

	return rv, nil
}

func (dw *DwtkAdapter) Close() error {
	return dw.device.Close()
}

func (dw *DwtkAdapter) Info() string {
	return fmt.Sprintf("Using dwtk custom hardware\nBaud Rate: %d bps\nBaud Rate Register: 0x%04x\n", dw.baudrate, dw.ubrr)
}
