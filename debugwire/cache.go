package debugwire

import (
	"fmt"
	"os"
)

type Cached struct {
	dw        *DebugWire
	Registers [32]byte
	PC        uint16
	SP        uint16
	SREG      byte
}

func (dw *DebugWire) Cache() (*Cached, error) {
	var err error
	rv := &Cached{dw: dw}

	rv.PC, err = dw.GetPC()
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nPC: 0x%04x\n\n", rv.PC)

	if err = dw.ReadRegisters(0, rv.Registers[:]); err != nil {
		return nil, err
	}

	// SPL, SPH, SREG
	sr := make([]byte, 3)
	if err = dw.ReadSRAM(0x5d, sr); err != nil {
		return nil, err
	}

	rv.SP = uint16(sr[1]<<8) | uint16(sr[0])
	rv.SREG = sr[2]

	return rv, nil
}

func (c *Cached) Restore() {
	rv := func() error {
		sr := []byte{
			byte(c.SP), byte(c.SP >> 8),
			c.SREG,
		}
		if err := c.dw.WriteSRAM(0x5d, sr); err != nil {
			return err
		}

		if err := c.dw.WriteRegisters(0, c.Registers[:]); err != nil {
			return err
		}

		return c.dw.SetPC(c.PC)
	}()

	if rv != nil {
		fmt.Fprintf(os.Stderr, "Error: debugwire: cache: %s\n", rv)
	}
}
