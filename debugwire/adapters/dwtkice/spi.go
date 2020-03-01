package dwtkice

import (
	"errors"
	"fmt"

	"github.com/dwtk/dwtk/avr"
)

type spiCommands struct {
	dev *device
}

func newSpiCommands(dev *device) *spiCommands {
	return &spiCommands{
		dev: dev,
	}
}

func (spi *spiCommands) enable() error {
	return spi.dev.controlIn(cmdSpiPgmEnable, 0, 0, nil)
}

func (spi *spiCommands) disable() error {
	return spi.dev.controlIn(cmdSpiPgmDisable, 0, 0, nil)
}

func (spi *spiCommands) command(c []byte) ([]byte, error) {
	if len(c) != 4 {
		return nil, fmt.Errorf("debugwire: dwtk-ice: invalid SPI command: %v", c)
	}
	f := make([]byte, 4)
	if err := spi.dev.controlIn(cmdSpiCommand, (uint16(c[0])<<8)|uint16(c[1]), (uint16(c[2])<<8)|uint16(c[3]), f); err != nil {
		return nil, err
	}
	return f, nil
}

func (spi *spiCommands) chipErase() error {
	_, err := spi.command(avr.SpiChipErase())
	return err
}

func (spi *spiCommands) readSignature() (uint16, error) {
	b, err := spi.command(avr.SpiReadSignature(0))
	if err != nil {
		return 0, err
	}
	if b[3] != 0x1e { // FIXME: move this value to avr package
		return 0, fmt.Errorf("debugwire: dwtk-ice: only devices manufactured by Atmel/Microchip are supported")
	}

	b, err = spi.command(avr.SpiReadSignature(1))
	if err != nil {
		return 0, err
	}
	rv := uint16(b[3]) << 8

	b, err = spi.command(avr.SpiReadSignature(2))
	if err != nil {
		return 0, err
	}
	return rv | uint16(b[3]), nil
}

func (spi *spiCommands) readLFuse() (byte, error) {
	b, err := spi.command(avr.SpiReadLFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (spi *spiCommands) writeLFuse(l byte) error {
	_, err := spi.command(avr.SpiWriteLFuse(l))
	return err
}

func (spi *spiCommands) readHFuse() (byte, error) {
	b, err := spi.command(avr.SpiReadHFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (spi *spiCommands) writeHFuse(l byte) error {
	_, err := spi.command(avr.SpiWriteHFuse(l))
	return err
}

func (spi *spiCommands) readEFuse() (byte, error) {
	b, err := spi.command(avr.SpiReadEFuse())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (spi *spiCommands) writeEFuse(l byte) error {
	_, err := spi.command(avr.SpiWriteEFuse(l))
	return err
}

func (spi *spiCommands) readLock() (byte, error) {
	b, err := spi.command(avr.SpiReadLock())
	if err != nil {
		return 0, err
	}
	return b[3], nil
}

func (spi *spiCommands) writeLock(l byte) error {
	_, err := spi.command(avr.SpiWriteLock(l))
	return err
}

func (spi *spiCommands) dwEnable(mcu *avr.MCU) error {
	if mcu == nil {
		return errors.New("debugwire: dwtk-ice: mcu not set")
	}
	f, err := spi.readHFuse()
	if err != nil {
		return err
	}
	f &= ^(mcu.DwenBit)
	return spi.writeHFuse(f)
}

func (spi *spiCommands) dwDisable(mcu *avr.MCU) error {
	if mcu == nil {
		return errors.New("debugwire: dwtk-ice: mcu not set")
	}
	if err := spi.enable(); err != nil {
		return err
	}
	f, err := spi.readHFuse()
	if err != nil {
		return err
	}
	f |= mcu.DwenBit
	return spi.writeHFuse(f)
}
