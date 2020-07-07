package usbserial

import (
	"fmt"
	"time"

	"github.com/dwtk/dwtk/avr"
)

func (us *UsbSerialAdapter) waitSPCMSR(mask byte) error {
	if err := us.SendBreak(); err != nil {
		return err
	}

	for i := 0; i < 0xff; i++ {
		if err := us.WriteInstruction(avr.IN(us.mcu.SPMCSR().Io8(), 29)); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond) // FIXME
		rv := make([]byte, 1)
		if err := us.ReadRegisters(29, rv); err != nil {
			return err
		}
		if rv[0]&mask == 0 {
			return nil
		}
	}

	return fmt.Errorf("debugwire: usbserial: flash: SPM timeout")
}

func (us *UsbSerialAdapter) spm() error {
	if err := us.WriteInstruction(avr.OUT(us.mcu.SPMCSR().Io8(), 29)); err != nil {
		return err
	}
	return us.WriteInstruction(avr.SPM())
}

func (us *UsbSerialAdapter) WriteFlashPage(start uint16, b []byte) error {
	c := []byte{
		avr.CTPB | avr.SPMEN,          // to set SPMCSR
		byte(start), byte(start >> 8), // Z
	}
	if err := us.WriteRegisters(29, c); err != nil {
		return err
	}
	if err := us.spm(); err != nil {
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
		if err := us.spm(); err != nil {
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
	if err := us.spm(); err != nil {
		return err
	}
	return us.waitSPCMSR(avr.SPMEN)
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
	if err := us.spm(); err != nil {
		return err
	}
	return us.waitSPCMSR(avr.SPMEN)
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
