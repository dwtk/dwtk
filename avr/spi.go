package avr

func SpiPgmEnable() []byte {
	return []byte{0xac, 0x53, 0x00, 0x00}
}

func SpiChipErase() []byte {
	return []byte{0xac, 0x80, 0x00, 0x00}
}

func SpiReadSignature(b byte) []byte {
	return []byte{0x30, 0x00, b, 0x00}
}

func SpiReadLFuse() []byte {
	return []byte{0x50, 0x00, 0x00, 0x00}
}

func SpiReadHFuse() []byte {
	return []byte{0x58, 0x08, 0x00, 0x00}
}

func SpiReadEFuse() []byte {
	return []byte{0x50, 0x08, 0x00, 0x00}
}

func SpiReadLock() []byte {
	return []byte{0x58, 0x00, 0x00, 0x00}
}

func SpiWriteLFuse(b byte) []byte {
	return []byte{0xac, 0xa0, 0x00, b}
}

func SpiWriteHFuse(b byte) []byte {
	return []byte{0xac, 0xa8, 0x00, b}
}

func SpiWriteEFuse(b byte) []byte {
	return []byte{0xac, 0xa4, 0x00, b}
}

func SpiWriteLock(b byte) []byte {
	return []byte{0xac, 0xe0, 0x00, b}
}
