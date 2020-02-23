package dwtkice

import (
	"errors"
	"fmt"
	"time"

	"github.com/dwtk/dwtk/avr"
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

	errNotSupportedDw  = errors.New("debugwire: dwtk-ice: operation not supported: target running on debugWIRE mode, try `dwtk disable`")
	errNotSupportedSpi = errors.New("debugwire: dwtk-ice: operation not supported: target running on SPI ISP mode, try `dwtk enable`")
)

type DwtkIceAdapter struct {
	device         *usbfs.Device
	ubrr           uint16
	targetBaudrate uint32
	actualBaudrate uint32
	spiMode        bool
}

func New() (*DwtkIceAdapter, error) {
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

	rv := &DwtkIceAdapter{
		device:  devices[0],
		spiMode: false,
	}
	if err := rv.device.Open(); err != nil {
		return nil, err
	}
	logger.Debug.Printf(" * Detected dwtk-ice %s", rv.device.GetVersion())

	if err := rv.controlIn(cmdDetectBaudrate, 0, 0, nil); err != nil {
		return nil, err
	}

	// we need a delay here to avoid issuing an usb request while dwtk-ice is
	// detecting baudrate with interrupts disabled.
	//
	// the math is easy:
	// - each bit takes 1/20Mhz seconds (50ns)
	// - max counter is 0xffff. 50ns * 0xffff ~= 3.3ms
	//
	// with some margin, we set 30ms because why not
	time.Sleep(30 * time.Millisecond)

	if errOrig := rv.controlGetError(); errOrig != nil {
		if err := rv.spiEnable(); err != nil {
			return nil, errOrig
		}
		rv.spiMode = true
		return rv, nil
	}

	f := make([]byte, 6)
	if err := rv.controlIn(cmdGetBaudrate, 0, 0, f); err != nil {
		return nil, err
	}

	if f[1] == 0 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid baudrate prescaler: 0")
	}
	if f[2] == 0 && f[3] == 0 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid pulse width: 0")
	}

	rv.ubrr = (uint16(f[4]) << 8) | uint16(f[5])
	rv.actualBaudrate = (uint32(f[0]) * 1000000) / uint32(uint16(f[1])*(rv.ubrr+1))
	rv.targetBaudrate = (uint32(f[0]) * 1000000) / uint32((uint16(f[2])<<8)|uint16(f[3]))

	logger.Debug.Printf(" * Actual baudrate: %d", rv.actualBaudrate)

	return rv, nil
}

func (dw *DwtkIceAdapter) Close() error {
	if dw.spiMode {
		dw.spiDisable()
	}
	return dw.device.Close()
}

func (dw *DwtkIceAdapter) Info() string {
	info := fmt.Sprintf("dwtk-ice %s\n", dw.device.GetVersion())
	if dw.spiMode {
		info += `
Target running on SPI ISP mode, try ` + "`dwtk enable`" + ` to enable debugWIRE mode.
`
	} else {
		info += fmt.Sprintf(`
Target running on debugWIRE mode, try `+"`dwtk disable`"+` to return to SPI ISP mode.

Target baudrate:   %d bps
Actual baudrate:   %d bps
Baudrate Register: 0x%04x
`,
			dw.targetBaudrate,
			dw.actualBaudrate,
			dw.ubrr)
	}
	return info
}

func (dw *DwtkIceAdapter) Enable() error {
	if !dw.spiMode {
		return errors.New("debugwire: dwtk-ice: target device is already running on debugWIRE mode")
	}

	if err := dw.spiEnable(); err != nil {
		return err
	}
	sign, err := dw.spiReadSignature()
	if err != nil {
		return err
	}
	mcu, err := avr.GetMCU(sign)
	if err != nil {
		return err
	}
	f, err := dw.spiReadHFuse()
	if err != nil {
		return err
	}
	f &= ^(mcu.DwenBit)
	if err := dw.spiWriteHFuse(f); err != nil {
		return err
	}

	fmt.Println("debugWIRE was enabled for target device. a target power cycle may be required")
	return nil
}

func (dw *DwtkIceAdapter) Disable() error {
	if dw.spiMode {
		return errors.New("debugwire: dwtk-ice: target device is already running on SPI ISP mode")
	}

	if err := dw.controlIn(cmdDisable, 0, 0, nil); err != nil {
		return err
	}
	// FIXME: check if SPI is supported
	if err := dw.spiEnable(); err != nil {
		return err
	}
	sign, err := dw.spiReadSignature()
	if err != nil {
		return err
	}
	mcu, err := avr.GetMCU(sign)
	if err != nil {
		return err
	}
	f, err := dw.spiReadHFuse()
	if err != nil {
		return err
	}
	f |= mcu.DwenBit
	if err := dw.spiWriteHFuse(f); err != nil {
		return err
	}

	fmt.Println("debugWIRE was disabled for target device. a target power cycle is required")
	return nil
}

func (dw *DwtkIceAdapter) Reset() error {
	if dw.spiMode {
		// FIXME: implement reset on firmware
		return nil
	}

	return dw.controlIn(cmdReset, 0, 0, nil)
}

func (dw *DwtkIceAdapter) ReadSignature() (uint16, error) {
	if dw.spiMode {
		return dw.spiReadSignature()
	}

	f := make([]byte, 2)
	if err := dw.controlIn(cmdReadSignature, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) ChipErase() error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spiChipErase()
}

func (dw *DwtkIceAdapter) SendBreak() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdSendBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) RecvBreak() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdRecvBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Go() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdGo, 0, 0, nil)
}

func (dw *DwtkIceAdapter) ResetAndGo() error {
	if err := dw.Reset(); err != nil {
		return err
	}
	if !dw.spiMode {
		return dw.Go()
	}
	return nil
}

func (dw *DwtkIceAdapter) Step() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdStep, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Continue(hwBreakpoint uint16, hwBreakpointSet bool, timers bool) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

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
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdWriteInstruction, inst, 0, nil)
}

func (dw *DwtkIceAdapter) WriteRegisters(start byte, regs []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlOut(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) ReadRegisters(start byte, regs []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) SetPC(pc uint16) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdSetPC, pc, 0, nil)
}

func (dw *DwtkIceAdapter) GetPC() (uint16, error) {
	if dw.spiMode {
		return 0, errNotSupportedSpi
	}

	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetPC, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) WriteSRAM(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlOut(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadSRAM(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) WriteFlashPage(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlOut(cmdWriteFlashPage, start, 0, data)
}

func (dw *DwtkIceAdapter) EraseFlashPage(start uint16) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdEraseFlashPage, start, 0, nil)
}

func (dw *DwtkIceAdapter) ReadFlash(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.controlIn(cmdReadFlash, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadFuses() ([]byte, error) {
	f := make([]byte, 4)
	if dw.spiMode {
		var err error
		f[0], err = dw.spiReadLFuse()
		if err != nil {
			return nil, err
		}
		f[1], err = dw.spiReadLock()
		if err != nil {
			return nil, err
		}
		f[2], err = dw.spiReadEFuse()
		if err != nil {
			return nil, err
		}
		f[3], err = dw.spiReadHFuse()
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if err := dw.controlIn(cmdReadFuses, 0, 0, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (dw *DwtkIceAdapter) WriteLFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spiWriteLFuse(data)
}

func (dw *DwtkIceAdapter) WriteHFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spiWriteHFuse(data)
}

func (dw *DwtkIceAdapter) WriteEFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spiWriteEFuse(data)
}

func (dw *DwtkIceAdapter) WriteLock(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spiWriteLock(data)
}
