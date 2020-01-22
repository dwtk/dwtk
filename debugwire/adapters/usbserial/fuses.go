package usbserial

import (
	"github.com/dwtk/dwtk/avr"
)

func (us *UsbSerialAdapter) ReadFuses() ([]byte, error) {
	b := []byte{
		avr.RFLB | avr.SPMEN, // to set SPMCSR
	}

	if err := us.WriteRegisters(30, []byte{0x00, 0x00}); err != nil {
		return nil, err
	}

	r := []byte{}
	d := make([]byte, 1)

	for i := 0; i < 4; i++ {
		if err := us.WriteRegisters(29, b); err != nil {
			return nil, err
		}

		if err := us.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
			return nil, err
		}

		if err := us.WriteInstruction(avr.LPM(28, true)); err != nil {
			return nil, err
		}

		if err := us.ReadRegisters(28, d); err != nil {
			return nil, err
		}

		r = append(r, d[0])
	}

	return r, nil
}
