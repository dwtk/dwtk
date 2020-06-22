package debugwire

import (
	"fmt"
	"os"
)

type regRange struct {
	start byte
	count byte
}

type cached struct {
	dw *DebugWIRE

	pc        uint16
	mask      uint32
	registers [32]byte
	ranges    []regRange
}

func (dw *DebugWIRE) cache(regs ...byte) (*cached, error) {
	if !dw.Cache {
		return &cached{
			dw: dw,
		}, nil
	}

	mask := uint32(0)
	for _, reg := range regs {
		if reg >= 32 {
			return nil, fmt.Errorf("debugwire: cache: invalid register: %d", reg)
		}
		mask |= (1 << reg)
	}

	pc, err := dw.adapter.GetPC()
	if err != nil {
		return nil, err
	}

	rv := &cached{
		dw:   dw,
		pc:   pc,
		mask: mask,
	}

	start := byte(0)
	count := byte(0)
	for i := byte(0); i < 32; i++ {
		if mask&(1<<i) != 0 {
			if count == 0 {
				start = i
			}
			count++
			if i == 31 {
				rv.ranges = append(rv.ranges, regRange{start, count})
			}
		} else if count > 0 {
			rv.ranges = append(rv.ranges, regRange{start, count})
			count = 0
			start = 0
		}
	}

	for _, r := range rv.ranges {
		if err := dw.adapter.ReadRegisters(r.start, rv.registers[r.start:r.start+r.count]); err != nil {
			return nil, err
		}
	}

	return rv, nil
}

func (c *cached) restore() {
	if !c.dw.Cache {
		return
	}

	rv := func() error {
		for _, r := range c.ranges {
			if err := c.dw.adapter.WriteRegisters(r.start, c.registers[r.start:r.start+r.count]); err != nil {
				return err
			}
		}

		return c.dw.adapter.SetPC(c.pc)
	}()

	if rv != nil {
		fmt.Fprintf(os.Stderr, "Error: debugwire: cache: %s\n", rv)
	}
}
