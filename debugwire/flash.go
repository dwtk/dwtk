package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWire) WriteFlashPage(start uint16, b []byte) error {
	if uint16(len(b)) != dw.MCU.FlashPageSize {
		return fmt.Errorf("debugwire: flash: page size must be 0x%04x for %s",
			dw.MCU.FlashPageSize,
			dw.MCU.Name,
		)
	}

	if start+dw.MCU.FlashPageSize > dw.MCU.FlashSize {
		return fmt.Errorf("debugwire: flash: writing out of flash space: 0x%04x + 0x%04x > 0x%04x",
			start,
			dw.MCU.FlashPageSize,
			dw.MCU.FlashSize,
		)
	}

	c := []byte{
		avr.CTPB | avr.SELFPRGEN,      // to set SPMCSR
		byte(start), byte(start >> 8), // Z
	}
	if err := dw.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.OUT(dw.MCU.SPMCSR, 29)); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	if err := dw.SendBreak(); err != nil {
		return err
	}

	c = []byte{
		avr.PGERS | avr.SELFPRGEN, // to set SPMCSR
	}
	if err := dw.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.OUT(dw.MCU.SPMCSR, 29)); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	if err := dw.SendBreak(); err != nil {
		return err
	}

	c = []byte{
		avr.SELFPRGEN, // to set SPMCSR
	}
	if err := dw.WriteRegisters(29, c); err != nil {
		return err
	}
	for i := 0; i < len(b); i += 2 {
		if err := dw.WriteRegisters(0, []byte{b[i], b[i+1]}); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.OUT(dw.MCU.SPMCSR, 29)); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.SPM()); err != nil {
			return err
		}
		if err := dw.WriteInstruction(avr.ADIW(30, 2)); err != nil {
			return err
		}
	}

	c = []byte{
		avr.PGWRT | avr.SELFPRGEN,     // to set SPMCSR
		byte(start), byte(start >> 8), // Z
	}
	if err := dw.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.OUT(dw.MCU.SPMCSR, 29)); err != nil {
		return err
	}
	if err := dw.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	return dw.SendBreak()
}

func (dw *DebugWire) WriteFlashInstruction(start uint16, inst uint16) error {
	c := []byte{
		byte(inst),
		byte(inst >> 8),
	}
	return dw.WriteFlash(start, c)
}

func (dw *DebugWire) WriteFlash(start uint16, b []byte) error {
	startPage := start / dw.MCU.FlashPageSize
	endAddr := start + uint16(len(b))
	endPage := (endAddr - 1) / dw.MCU.FlashPageSize

	k := 0
	for i := startPage; i <= endPage; i += 1 {
		addr := i * dw.MCU.FlashPageSize
		page := make([]byte, dw.MCU.FlashPageSize)
		if err := dw.ReadFlash(addr, page); err != nil {
			return err
		}
		pStart := uint16(0)
		if start >= addr {
			pStart = start - addr
		}
		pEnd := endAddr - addr
		if pEnd > dw.MCU.FlashPageSize {
			pEnd = dw.MCU.FlashPageSize
		}

		for j := pStart; j < pEnd; j += 1 {
			page[j] = b[k]
			k += 1
		}

		if err := dw.WriteFlashPage(addr, page); err != nil {
			return err
		}
	}

	return nil
}

func (dw *DebugWire) ReadFlash(start uint16, b []byte) error {
	if start+uint16(len(b)) > dw.MCU.FlashSize {
		return fmt.Errorf("debugwire: flash: reading out of flash space: 0x%04x + 0x%04x > 0x%04x",
			start,
			len(b),
			dw.MCU.FlashSize,
		)
	}

	c := []byte{
		byte(start), byte(start >> 8), // Z
	}
	if err := dw.WriteRegisters(30, c); err != nil {
		return err
	}

	l := uint16(len(b) * 2)
	c = []byte{
		0x66,
		0xd0, 0x00, 0x00,
		0xd1, byte(l >> 8), byte(l),
		0xc2, 0x02,
		0x20,
	}
	if err := dw.Port.Write(c); err != nil {
		return err
	}

	return dw.Port.Read(b)
}
