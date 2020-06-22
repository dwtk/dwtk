package debugwire

import (
	"fmt"
	"time"

	"github.com/dwtk/dwtk/avr"
)

func (dw *DebugWIRE) readEEPROM(start uint16, b []byte) error {
	c := []byte{
		avr.EERE,
		byte(start), byte(start >> 8),
	}
	if err := dw.adapter.WriteRegisters(29, c); err != nil {
		return nil
	}

	d := make([]byte, 1)
	for i := 0; i < len(b); i++ {
		// EEARL
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR+2, 30)); err != nil {
			return err
		}
		if dw.MCU.WithEEARH {
			// EEARH
			if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR+3, 31)); err != nil {
				return err
			}
		}
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR, 29)); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.ADIW(30, 1)); err != nil {
			return err
		}
		// EEDR
		if err := dw.WriteInstruction(avr.IN(dw.MCU.EECR+1, 0)); err != nil {
			return err
		}
		if err := dw.adapter.ReadRegisters(0, d); err != nil {
			return err
		}
		b[i] = d[0]
	}

	return nil
}

func (dw *DebugWIRE) WriteEEPROM(start uint16, b []byte) error {
	if start+uint16(len(b)) > dw.MCU.EEPROMSize {
		return fmt.Errorf("debugwire: eeprom: writing out of eeprom space: 0x%04x + 0x%04x > 0x%04x",
			start,
			len(b),
			dw.MCU.EEPROMSize,
		)
	}

	cache, err := dw.cache(0, 1, 28, 29, 30, 31)
	if err != nil {
		return err
	}
	defer cache.restore()

	w := make([]byte, len(b))
	if err := dw.readEEPROM(start, w); err != nil {
		return err
	}

	c := []byte{
		avr.EEMPE,
		avr.EEPE,
		byte(start), byte(start >> 8),
	}
	if err := dw.adapter.WriteRegisters(28, c); err != nil {
		return nil
	}

	for i := 0; i < len(b); i++ {
		if b[i] == w[i] { // do not write unless needed
			if err := dw.WriteInstruction(avr.ADIW(30, 1)); err != nil {
				return err
			}
			continue
		}

		// EEARL
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR+2, 30)); err != nil {
			return err
		}
		//
		if dw.MCU.WithEEARH {
			// EEARH
			if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR+3, 31)); err != nil {
				return err
			}
		}
		if err := dw.adapter.WriteRegisters(0, []byte{b[i]}); err != nil {
			return nil
		}
		// EEDR
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR+1, 0)); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.ADIW(30, 1)); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR, 28)); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.EECR, 29)); err != nil {
			return err
		}
		if err := dw.SendBreak(); err != nil {
			return err
		}

		time.Sleep(5 * time.Millisecond)
	}

	return nil
}

func (dw *DebugWIRE) ReadEEPROM(start uint16, b []byte) error {
	if start+uint16(len(b)) > dw.MCU.EEPROMSize {
		return fmt.Errorf("debugwire: eeprom: reading out of eeprom space: 0x%04x + 0x%04x > 0x%04x",
			start,
			len(b),
			dw.MCU.EEPROMSize,
		)
	}

	cache, err := dw.cache(0, 1, 29, 30, 31)
	if err != nil {
		return err
	}
	defer cache.restore()

	return dw.readEEPROM(start, b)
}
