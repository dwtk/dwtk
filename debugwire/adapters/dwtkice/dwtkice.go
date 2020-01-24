package dwtkice

import (
	"errors"
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

	iceErrors = map[uint8]error{
		1: errors.New("debugwire: dwtk-ice: baudrate detection failed"),
		2: errors.New("debugwire: dwtk-ice: target detection failed"),
		3: errors.New("debugwire: dwtk-ice: got unexpected target device"),
		4: errors.New("debugwire: dwtk-ice: got unexpected byte echoed back"),
		5: errors.New("debugwire: dwtk-ice: got unexpected break value"),
		6: errors.New("debugwire: dwtk-ice: read/write data is too large"),
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

func (dw *DwtkIceAdapter) Disable() error {
	return dw.controlIn(cmdDisable, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Reset() error {
	return dw.controlIn(cmdReset, 0, 0, nil)
}

func (dw *DwtkIceAdapter) GetSignature() (uint16, error) {
	return dw.signature, nil
}

func (dw *DwtkIceAdapter) SendBreak() error {
	return dw.controlIn(cmdSendBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) RecvBreak() error {
	return dw.controlIn(cmdRecvBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Go() error {
	return dw.controlIn(cmdGo, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Step() error {
	return dw.controlIn(cmdStep, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Continue(hwBreakpoint uint16, hwBreakpointSet bool, timers bool) error {
	// idx: byte 0 -> hw bp set
	//      byte 1 -> timers
	idx := uint16(0)
	if hwBreakpointSet {
		idx |= (1 << 0)
	}
	if timers {
		idx |= (1 << 1)
	}
	return dw.controlIn(cmdContinue, hwBreakpoint, idx, nil)
}

func (dw *DwtkIceAdapter) WriteInstruction(inst uint16) error {
	return dw.controlIn(cmdWriteInstruction, inst, 0, nil)
}

func (dw *DwtkIceAdapter) WriteRegisters(start byte, regs []byte) error {
	return dw.controlOut(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) ReadRegisters(start byte, regs []byte) error {
	return dw.controlIn(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) SetPC(pc uint16) error {
	return dw.controlIn(cmdSetPC, pc, 0, nil)
}

func (dw *DwtkIceAdapter) GetPC() (uint16, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetPC, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) WriteSRAM(start uint16, data []byte) error {
	return dw.controlOut(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadSRAM(start uint16, data []byte) error {
	return dw.controlIn(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) WriteFlashPage(start uint16, data []byte) error {
	return dw.controlOut(cmdWriteFlashPage, start, 0, data)
}

func (dw *DwtkIceAdapter) EraseFlashPage(start uint16) error {
	return dw.controlIn(cmdEraseFlashPage, start, 0, nil)
}

func (dw *DwtkIceAdapter) ReadFlash(start uint16, data []byte) error {
	return dw.controlIn(cmdReadFlash, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadFuses() ([]byte, error) {
	f := make([]byte, 4)
	if err := dw.controlIn(cmdReadFuses, 0, 0, f); err != nil {
		return nil, err
	}
	return f, nil
}
