package dwtkice

import (
	"fmt"

	"github.com/dwtk/dwtk/avr"
)

func (dw *DwtkIceAdapter) spiEnable() error {
	return dw.controlIn(cmdSpiPgmEnable, 0, 0, nil)
}

func (dw *DwtkIceAdapter) spiDisable() error {
	return dw.controlIn(cmdSpiPgmDisable, 0, 0, nil)
}

func (dw *DwtkIceAdapter) spiCommand(c []byte) ([]byte, error) {
	if len(c) != 4 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid SPI command: %v", c)
	}
	f := make([]byte, 4)
	if err := dw.controlIn(cmdSpiCommand, (uint16(c[0])<<8)|uint16(c[1]), (uint16(c[2])<<8)|uint16(c[3]), f); err != nil {
		return nil, err
	}
	return f, nil
}

func (dw *DwtkIceAdapter) spiChipErase() error {
	_, err := dw.spiCommand(avr.SpiChipErase())
	return err
}

func (dw *DwtkIceAdapter) spiReadSignature() (uint16, error) {
	b, err := dw.spiCommand(avr.SpiReadSignature(0))
	if err != nil {
		return 0, err
	}
	if b[3] != 0x1e { // FIXME: move this value to avr package
		return 0, fmt.Errorf("debugwire: dwtk-ice: only devices manufactured by Atmel/Microchip are supported")
	}

	b, err = dw.spiCommand(avr.SpiReadSignature(1))
	if err != nil {
		return 0, err
	}
	rv := uint16(b[3]) << 8

	b, err = dw.spiCommand(avr.SpiReadSignature(2))
	if err != nil {
		return 0, err
	}
	return rv | uint16(b[3]), nil
}

func (dw *DwtkIceAdapter) spiReadLFuse() (byte, error) {
	b, err := dw.spiCommand(avr.SpiReadLFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (dw *DwtkIceAdapter) spiWriteLFuse(l byte) error {
	_, err := dw.spiCommand(avr.SpiWriteLFuse(l))
	return err
}

func (dw *DwtkIceAdapter) spiReadHFuse() (byte, error) {
	b, err := dw.spiCommand(avr.SpiReadHFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (dw *DwtkIceAdapter) spiWriteHFuse(l byte) error {
	_, err := dw.spiCommand(avr.SpiWriteHFuse(l))
	return err
}

func (dw *DwtkIceAdapter) spiReadEFuse() (byte, error) {
	b, err := dw.spiCommand(avr.SpiReadEFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (dw *DwtkIceAdapter) spiWriteEFuse(l byte) error {
	_, err := dw.spiCommand(avr.SpiWriteEFuse(l))
	return err
}

func (dw *DwtkIceAdapter) spiReadLock() (byte, error) {
	b, err := dw.spiCommand(avr.SpiReadLock())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (dw *DwtkIceAdapter) spiWriteLock(l byte) error {
	_, err := dw.spiCommand(avr.SpiWriteLock(l))
	return err
}

func (dw *DwtkIceAdapter) spiEnableDw(mcu *avr.MCU) error {
	f, err := dw.spiReadHFuse()
	if err != nil {
		return err
	}
	f &= ^(mcu.DwenBit)
	return dw.spiWriteHFuse(f)
}

func (dw *DwtkIceAdapter) spiDisableDw(mcu *avr.MCU) error {
	f, err := dw.spiReadHFuse()
	if err != nil {
		return err
	}
	f |= mcu.DwenBit
	return dw.spiWriteHFuse(f)
}
