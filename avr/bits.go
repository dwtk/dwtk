package avr

const (
	SPMEN  = byte(1 << 0)
	PGERS  = byte(1 << 1)
	PGWRT  = byte(1 << 2)
	RFLB   = byte(1 << 3)
	CTPB   = byte(1 << 4)
	RWWSRE = byte(1 << 4)
	RWWSB  = byte(1 << 6)
)

const (
	EERE  = byte(1 << 0)
	EEPE  = byte(1 << 1)
	EEMPE = byte(1 << 2)
)

const (
	LOW_FUSE      = 0x00
	LOCKBIT       = 0x01
	EXTENDED_FUSE = 0x02
	HIGH_FUSE     = 0x03
)
