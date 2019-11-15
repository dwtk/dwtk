package avr

var (
	// SPMCSR
	SELFPRGEN = byte(1 << 0)
	PGERS     = byte(1 << 1)
	PGWRT     = byte(1 << 2)
	RFLB      = byte(1 << 3)
	CTPB      = byte(1 << 4)
)
