package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWIRE) WriteFlashPage(start uint16, b []byte) error {
	if uint16(len(b)) != dw.MCU.FlashPageSize {
		return fmt.Errorf("debugwire: flash: page size must be 0x%04x for %s",
			dw.MCU.FlashPageSize,
			dw.MCU.Name,
		)
	}

	if start%dw.MCU.FlashPageSize != 0 {
		return fmt.Errorf("debugwire: flash: start address must be aligned to page start (page size: 0x%04x)",
			dw.MCU.FlashPageSize,
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
		avr.CTPB | avr.SPMEN,          // to set SPMCSR
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
		avr.PGERS | avr.SPMEN, // to set SPMCSR
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
		avr.SPMEN, // to set SPMCSR
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
		avr.PGWRT | avr.SPMEN,         // to set SPMCSR
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

func (dw *DebugWIRE) WriteFlashInstruction(start uint16, inst uint16) error {
	c := []byte{
		byte(inst),
		byte(inst >> 8),
	}
	return dw.WriteFlash(start, c)
}

func (dw *DebugWIRE) WriteFlash(start uint16, b []byte) error {
	startPage := start / dw.MCU.FlashPageSize
	endAddr := start + uint16(len(b))
	endPage := (endAddr - 1) / dw.MCU.FlashPageSize

	pages := make(map[int][]byte)

	for i := startPage; i <= endPage; i++ {
		addr := i * dw.MCU.FlashPageSize
		page := make([]byte, dw.MCU.FlashPageSize)

		if err := dw.ReadFlash(addr, page); err != nil {
			return err
		}

		pages[int(i)] = page
	}

	k := 0
	for i := startPage; i <= endPage; i++ {
		addr := i * dw.MCU.FlashPageSize
		page, ok := pages[int(i)]
		if !ok {
			return fmt.Errorf("debugwire: flash: bad page split")
		}

		pStart := uint16(0)
		if start >= addr {
			pStart = start - addr
		}
		pEnd := endAddr - addr
		if pEnd > dw.MCU.FlashPageSize {
			pEnd = dw.MCU.FlashPageSize
		}

		for j := pStart; j < pEnd; j++ {
			page[j] = b[k]
			k++
		}

		if err := dw.WriteFlashPage(addr, page); err != nil {
			return err
		}
	}

	return nil
}

func (dw *DebugWIRE) ReadFlash(start uint16, b []byte) error {
	if start+uint16(len(b)) > dw.MCU.FlashSize {
		return fmt.Errorf("debugwire: flash: reading out of flash space: 0x%04x + 0x%04x > 0x%04x",
			start,
			len(b),
			dw.MCU.FlashSize,
		)
	}

	return dw.adapter.ReadFlash(start, b)
}
