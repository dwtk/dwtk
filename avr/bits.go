package avr

const (
	// SPMCSR
	SELFPRGEN = byte(1 << 0)
	PGERS     = byte(1 << 1)
	PGWRT     = byte(1 << 2)
	RFLB      = byte(1 << 3)
	CTPB      = byte(1 << 4)
)

const (
	SPL  = 0x5d
	SPH  = 0x5e
	SREG = 0x5f
)
