package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWire) WriteFlashPage(start uint16, b []byte) error {
	if uint16(len(b)) != dw.MCU.FlashPageSize {
		return fmt.Errorf("debugwire: flash page size, must be 0x%02x for %s",
			dw.MCU.FlashPageSize,
			dw.MCU.Name,
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

func (dw *DebugWire) ReadFlash(start uint16, b []byte) error {
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
