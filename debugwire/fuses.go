package debugwire

import (
	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWIRE) ReadFuses() ([]byte, error) {
	b := []byte{
		avr.RFLB | avr.SELFPRGEN, // to set SPMCSR
		0x00, 0x00,               // Z
	}

	f := []byte{
		0x00, // low
		0x03, // high
		0x02, // extended
		0x01, // lockbit
	}

	r := []byte{}

	for _, i := range f {
		b[1] = i // Z lo

		if err := dw.WriteRegisters(29, b); err != nil {
			return nil, err
		}

		if err := dw.WriteInstruction(avr.OUT(dw.MCU.SPMCSR, 29)); err != nil {
			return nil, err
		}

		if err := dw.WriteInstruction(avr.LPM(0)); err != nil {
			return nil, err
		}

		d := make([]byte, 1)
		if err := dw.ReadRegisters(0, d); err != nil {
			return nil, err
		}

		r = append(r, d[0])
	}

	return r, nil
}
