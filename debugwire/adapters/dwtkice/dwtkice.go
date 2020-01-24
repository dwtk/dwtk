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
	cmdGetError = iota + 1
	cmdInit
	cmdDisable
	cmdReset
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
		cmdGetError:         "cmdGetError",
		cmdInit:             "cmdInit",
		cmdDisable:          "cmdDisable",
		cmdReset:            "cmdReset",
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
		cmdReadFuses:        "cmdReadFuses",
	}
)

type DwtkIceAdapter struct {
	device         *usbfs.Device
	ubrr           uint16
	targetBaudrate uint32
	actualBaudrate uint32
	signature      uint16
}

func New() (*DwtkIceAdapter, error) {
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
		device: devices[0],
	}
	if err := rv.device.Open(); err != nil {
		return nil, err
	}
	logger.Debug.Printf(" * Detected dwtk-ice %s", rv.device.GetVersion())

	f := make([]byte, 8)
	if err := rv.controlIn(cmdInit, 0, 0, f); err != nil {
		return nil, err
	}

	rv.ubrr = (uint16(f[4]) << 8) | uint16(f[5])
	rv.actualBaudrate = (uint32(f[0]) * 1000000) / uint32(uint16(f[1])*(rv.ubrr+1))
	rv.targetBaudrate = (uint32(f[0]) * 1000000) / uint32((uint16(f[2])<<8)|uint16(f[3]))
	rv.signature = (uint16(f[6]) << 8) | uint16(f[7])

	logger.Debug.Printf(" * Actual baudrate: %d", rv.actualBaudrate)

	return rv, nil
}

func (dw *DwtkIceAdapter) Close() error {
	return dw.device.Close()
}

func (dw *DwtkIceAdapter) Info() string {
	return fmt.Sprintf(
		`dwtk-ice %s

Target baudrate:   %d bps
Actual baudrate:   %d bps
Baudrate Register: 0x%04x
`,
		dw.device.GetVersion(),
		dw.targetBaudrate,
		dw.actualBaudrate,
		dw.ubrr,
	)
}

func (dw *DwtkIceAdapter) GetSignature() (uint16, error) {
	return dw.signature, nil
}
