package debugwire

import (
	"golang.rgm.io/dwtk/avr"
)

func (dw *DebugWIRE) ReadFuses() ([]byte, error) {
	b := []byte{
		avr.RFLB | avr.SPMEN, // to set SPMCSR
	}

	if err := dw.WriteRegisters(29, append(b, 0x00, 0x00)); err != nil {
		return nil, err
	}

	r := []byte{}
	d := make([]byte, 1)

	for i := 0; i < 4; i++ {
		if i > 0 {
			if err := dw.WriteRegisters(29, b); err != nil {
				return nil, err
			}
		}

		if err := dw.WriteInstruction(avr.OUT(avr.SPMCSR, 29)); err != nil {
			return nil, err
		}

		if err := dw.WriteInstruction(avr.LPM(28, true)); err != nil {
			return nil, err
		}

		if err := dw.ReadRegisters(28, d); err != nil {
			return nil, err
		}

		r = append(r, d[0])
	}

	return r, nil
}
