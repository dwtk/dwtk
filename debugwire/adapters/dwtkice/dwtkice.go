package dwtkice

import (
	"errors"
	"fmt"
	"time"

	"github.com/dwtk/dwtk/avr"
	"github.com/dwtk/dwtk/internal/logger"
)

var (
	errNotSupportedDw  = errors.New("debugwire: dwtk-ice: operation not supported: target running on debugWIRE mode, try `dwtk disable`")
	errNotSupportedSpi = errors.New("debugwire: dwtk-ice: operation not supported: target running on SPI ISP mode, try `dwtk enable`")
)

type DwtkIceAdapter struct {
	dev            *device
	spi            *spiCommands
	mcu            *avr.MCU
	ubrr           uint16
	targetBaudrate uint32
	actualBaudrate uint32
	spiMode        bool
}

func New() (*DwtkIceAdapter, error) {
	dev, err := newDevice()
	if dev == nil || err != nil {
		return nil, err
	}
	serial := dev.getSerial()
	if serial != "" {
		logger.Debug.Printf(" * Detected dwtk-ice %s (SN: %s)", dev.getVersion(), serial)
	} else {
		logger.Debug.Printf(" * Detected dwtk-ice %s", dev.getVersion())
	}

	var spi *spiCommands
	if dev.spi {
		spi = newSpiCommands(dev)
	}

	if err := dev.controlIn(cmdDetectBaudrate, 0, 0, nil); err != nil {
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

	if errOrig := dev.controlGetError(); errOrig != nil {
		if !dev.spi {
			return nil, errOrig
		}
		if err := spi.enable(); err != nil {
			return nil, errOrig
		}
		return &DwtkIceAdapter{dev: dev, spi: spi, spiMode: true}, nil
	}

	f := make([]byte, 6)
	if err := dev.controlIn(cmdGetBaudrate, 0, 0, f); err != nil {
		return nil, err
	}

	if f[1] == 0 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid baudrate prescaler: 0")
	}
	if f[2] == 0 && f[3] == 0 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid pulse width: 0")
	}

	ubrr := (uint16(f[4]) << 8) | uint16(f[5])
	rv := &DwtkIceAdapter{
		dev:            dev,
		spi:            spi,
		ubrr:           ubrr,
		actualBaudrate: (uint32(f[0]) * 1000000) / uint32(uint16(f[1])*(ubrr+1)),
		targetBaudrate: (uint32(f[0]) * 1000000) / uint32((uint16(f[2])<<8)|uint16(f[3])),
		spiMode:        false,
	}

	logger.Debug.Printf(" * Actual baudrate: %d", rv.actualBaudrate)
	return rv, nil
}

func (dw *DwtkIceAdapter) Close() error {
	if dw.spiMode {
		dw.spi.disable()
	}
	return dw.dev.close()
}

func (dw *DwtkIceAdapter) Info() string {
	info := ""
	serial := dw.dev.getSerial()
	if serial != "" {
		info = fmt.Sprintf("dwtk-ice %s (SN: %s)\n", dw.dev.getVersion(), dw.dev.getSerial())
	} else {
		info = fmt.Sprintf("dwtk-ice %s\n", dw.dev.getVersion())
	}
	if dw.spiMode {
		info += "\nTarget running on SPI ISP mode, try `dwtk enable` to enable debugWIRE mode.\n"
	} else {
		if dw.dev.spi {
			info += "\nTarget running on debugWIRE mode, try `dwtk disable` to return to SPI ISP mode.\n"
		}
		info += fmt.Sprintf(`
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

func (dw *DwtkIceAdapter) SetMCU(mcu *avr.MCU) {
	dw.mcu = mcu
}

func (dw *DwtkIceAdapter) Enable() error {
	if !dw.spiMode {
		if err := dw.ResetAndGo(); err != nil {
			return err
		}
		return errors.New("debugwire: dwtk-ice: target device is already running on debugWIRE mode")
	}
	if err := dw.spi.dwEnable(dw.mcu); err != nil {
		return err
	}
	if err := dw.dev.controlIn(cmdSpiReset, 0, 0, nil); err != nil {
		return err
	}

	fmt.Println("debugWIRE was enabled for target device. a target power cycle may be required")
	return nil
}

func (dw *DwtkIceAdapter) Disable() error {
	if dw.spiMode {
		if err := dw.ResetAndGo(); err != nil {
			return err
		}
		return errors.New("debugwire: dwtk-ice: target device is already running on SPI ISP mode")
	}
	if err := dw.dev.controlIn(cmdDisable, 0, 0, nil); err != nil {
		return err
	}
	if dw.dev.spi {
		if err := dw.spi.dwDisable(dw.mcu); err == nil {
			if err := dw.dev.controlIn(cmdSpiReset, 0, 0, nil); err != nil {
				return err
			}
			fmt.Println("debugWIRE was disabled for target device. a target power cycle is required")
			return nil
		}
	}

	fmt.Println("debugWIRE was disabled for target device, and it can be flashed using an SPI ISP now.")
	fmt.Println("this must be done without a target power cycle.")
	return nil
}

func (dw *DwtkIceAdapter) Reset() error {
	if dw.spiMode {
		return dw.dev.controlIn(cmdSpiReset, 0, 0, nil)
	}

	return dw.dev.controlIn(cmdReset, 0, 0, nil)
}

func (dw *DwtkIceAdapter) ReadSignature() (uint16, error) {
	if dw.spiMode {
		return dw.spi.readSignature()
	}

	f := make([]byte, 2)
	if err := dw.dev.controlIn(cmdReadSignature, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) ChipErase() error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spi.chipErase()
}

func (dw *DwtkIceAdapter) SendBreak() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdSendBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) RecvBreak() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdRecvBreak, 0, 0, nil)
}

func (dw *DwtkIceAdapter) Go() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdGo, 0, 0, nil)
}

func (dw *DwtkIceAdapter) ResetAndGo() error {
	if err := dw.Reset(); err != nil {
		return err
	}
	if !dw.spiMode {
		// we don't have a go command for spi, afaik, resetting starts program right away.
		// i won't make the go command return nil, because the command itself doesn't exists for spi
		return dw.Go()
	}
	return nil
}

func (dw *DwtkIceAdapter) Step() error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdStep, 0, 0, nil)
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
	return dw.dev.controlIn(cmdContinue, hwBreakpoint, idx, nil)
}

func (dw *DwtkIceAdapter) WriteInstruction(inst uint16) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdWriteInstruction, inst, 0, nil)
}

func (dw *DwtkIceAdapter) WriteRegisters(start byte, regs []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlOut(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) ReadRegisters(start byte, regs []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdRegisters, uint16(start), 0, regs)
}

func (dw *DwtkIceAdapter) SetPC(pc uint16) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdSetPC, pc, 0, nil)
}

func (dw *DwtkIceAdapter) GetPC() (uint16, error) {
	if dw.spiMode {
		return 0, errNotSupportedSpi
	}

	f := make([]byte, 2)
	if err := dw.dev.controlIn(cmdGetPC, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) WriteSRAM(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlOut(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadSRAM(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdSRAM, start, 0, data)
}

func (dw *DwtkIceAdapter) WriteFlashPage(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlOut(cmdWriteFlashPage, start, 0, data)
}

func (dw *DwtkIceAdapter) EraseFlashPage(start uint16) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdEraseFlashPage, start, 0, nil)
}

func (dw *DwtkIceAdapter) ReadFlash(start uint16, data []byte) error {
	if dw.spiMode {
		return errNotSupportedSpi
	}

	return dw.dev.controlIn(cmdReadFlash, start, 0, data)
}

func (dw *DwtkIceAdapter) ReadFuses() ([]byte, error) {
	f := make([]byte, 4)
	if dw.spiMode {
		var err error
		f[0], err = dw.spi.readLFuse()
		if err != nil {
			return nil, err
		}
		f[1], err = dw.spi.readLock()
		if err != nil {
			return nil, err
		}
		f[2], err = dw.spi.readEFuse()
		if err != nil {
			return nil, err
		}
		f[3], err = dw.spi.readHFuse()
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if err := dw.dev.controlIn(cmdReadFuses, 0, 0, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (dw *DwtkIceAdapter) WriteLFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spi.writeLFuse(data)
}

func (dw *DwtkIceAdapter) WriteHFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spi.writeHFuse(data)
}

func (dw *DwtkIceAdapter) WriteEFuse(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spi.writeEFuse(data)
}

func (dw *DwtkIceAdapter) WriteLock(data byte) error {
	if !dw.spiMode {
		return errNotSupportedDw
	}

	return dw.spi.writeLock(data)
}
