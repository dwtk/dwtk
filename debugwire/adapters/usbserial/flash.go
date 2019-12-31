package usbserial

import (
	"golang.rgm.io/dwtk/avr"
)

func (us *UsbSerialAdapter) WriteFlashPage(start uint16, b []byte) error {
	c := []byte{
		avr.CTPB | avr.SPMEN,          // to set SPMCSR
		byte(start), byte(start >> 8), // Z
	}
	if err := us.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	if err := us.SendBreak(); err != nil {
		return err
	}

	if err := us.eraseFlashPage(0, false); err != nil {
		return err
	}

	c = []byte{
		avr.SPMEN, // to set SPMCSR
	}
	if err := us.WriteRegisters(29, c); err != nil {
		return err
	}
	for i := 0; i < len(b); i += 2 {
		if err := us.WriteRegisters(0, []byte{b[i], b[i+1]}); err != nil {
			return err
		}
		if err := us.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
			return err
		}
		if err := us.WriteInstruction(avr.SPM()); err != nil {
			return err
		}
		if err := us.WriteInstruction(avr.ADIW(30, 2)); err != nil {
			return err
		}
	}

	c = []byte{
		avr.PGWRT | avr.SPMEN,         // to set SPMCSR
		byte(start), byte(start >> 8), // Z
	}
	if err := us.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	return us.SendBreak()
}

func (us *UsbSerialAdapter) eraseFlashPage(start uint16, setStart bool) error {
	c := []byte{
		avr.PGERS | avr.SPMEN, // to set SPMCSR
	}
	if setStart {
		c = append(c, byte(start), byte(start>>8))
	}
	if err := us.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
		return err
	}
	if err := us.WriteInstruction(avr.SPM()); err != nil {
		return err
	}
	if err := us.SendBreak(); err != nil {
		return err
	}
	return nil
}

func (us *UsbSerialAdapter) EraseFlashPage(start uint16) error {
	return us.eraseFlashPage(start, true)
}

func (us *UsbSerialAdapter) ReadFlash(start uint16, b []byte) error {
	c := []byte{
		byte(start), byte(start >> 8), // Z
	}
	if err := us.WriteRegisters(30, c); err != nil {
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
	if err := us.device.Write(c); err != nil {
		return err
	}

	return us.device.Read(b)
}
