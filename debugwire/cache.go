package debugwire

import (
	"fmt"
	"os"
)

type Cached struct {
	R0   byte
	R1   byte
	R29  byte
	R30  byte
	R31  byte
	PC   uint16
	SREG byte

	dw     *DebugWire
	with01 bool
}

func (dw *DebugWire) Cache(with01 bool) (*Cached, error) {
	var err error
	rv := &Cached{
		dw:     dw,
		with01: with01,
	}

	rv.PC, err = dw.GetPC()
	if err != nil {
		return nil, err
	}

	if with01 {
		b := make([]byte, 2)
		if err := dw.ReadRegisters(0, b); err != nil {
			return nil, err
		}
		rv.R0 = b[0]
		rv.R1 = b[1]
	}

	b := make([]byte, 3)
	if err = dw.ReadRegisters(29, b); err != nil {
		return nil, err
	}
	rv.R29 = b[0]
	rv.R30 = b[1]
	rv.R31 = b[2]

	rv.SREG, err = dw.GetSREG()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (c *Cached) Restore() {
	rv := func() error {
		if err := c.dw.SetSREG(c.SREG); err != nil {
			return err
		}

		if c.with01 {
			if err := c.dw.WriteRegisters(0, []byte{c.R0, c.R1}); err != nil {
				return err
			}
		}

		if err := c.dw.WriteRegisters(29, []byte{c.R29, c.R30, c.R31}); err != nil {
			return err
		}

		return c.dw.SetPC(c.PC)
	}()

	if rv != nil {
		fmt.Fprintf(os.Stderr, "Error: debugwire: cache: %s\n", rv)
	}
}
