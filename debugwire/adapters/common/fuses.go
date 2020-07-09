package common

import (
	"errors"

	"github.com/dwtk/dwtk/avr"
)

func ReadFuses(c Common) ([]byte, error) {
	b := []byte{
		avr.RFLB | avr.SPMEN, // to set SPMCSR
	}

	if err := c.WriteRegisters(30, []byte{0x00, 0x00}); err != nil {
		return nil, err
	}

	r := []byte{}
	d := make([]byte, 1)

	for i := 0; i < 4; i++ {
		if err := c.WriteRegisters(29, b); err != nil {
			return nil, err
		}

		mcu := c.GetMCU()
		if mcu == nil {
			return nil, errors.New("debugwire: MCU not set")
		}

		if err := c.WriteInstruction(avr.OUT(mcu.SPMCSR().Io8(), 29)); err != nil {
			return nil, err
		}

		if err := c.WriteInstruction(avr.LPM(28, true)); err != nil {
			return nil, err
		}

		if err := c.ReadRegisters(28, d); err != nil {
			return nil, err
		}

		r = append(r, d[0])
	}

	return r, nil
}
