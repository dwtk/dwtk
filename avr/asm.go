package avr

func ADIW(reg byte, val uint16) uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_ADIW.html
	// opcode: 1001 0110 KKdd KKKK
	op := uint16(0b1001011000000000)
	kh := uint16(0b110000 & val)
	kh <<= 2
	kl := uint16(0b1111 & val)
	reg -= byte(24)
	reg >>= 1
	de := uint16(0b11 & reg)
	de <<= 4
	return op | kh | kl | de
}

func BREAK() uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_BREAK.html
	// 1001 0101 1001 1000
	op := uint16(0b1001010110011000)
	return op
}

func IN(addr byte, reg byte) uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_IN.html
	// opcode: 1011 0AAd dddd AAAA
	op := uint16(0b1011000000000000)
	ah := uint16(0b110000 & addr)
	ah <<= 5
	al := uint16(0b1111 & addr)
	de := uint16(0b11111 & reg)
	de <<= 4
	return op | ah | al | de
}

func OUT(addr byte, reg byte) uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_OUT.html
	// opcode: 1011 1AAr rrrr AAAA
	op := uint16(0b1011100000000000)
	ah := uint16(0b110000 & addr)
	ah <<= 5
	al := uint16(0b1111 & addr)
	re := uint16(0b11111 & reg)
	re <<= 4
	return op | ah | al | re
}

func LPM(reg byte, incr bool) uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_LPM.html
	// opcode: 1001 000d dddd 0100 - Z
	// opcode: 1001 000d dddd 0101 - Z+
	op := uint16(0b1001000000000100)
	if incr {
		op |= 0b1
	}
	de := uint16(0b11111 & reg)
	de <<= 4
	return op | de
}

func SPM() uint16 {
	// https://www.microchip.com/webdoc/avrassembler/avrassembler.wb_SPM.html
	// opcode: 1001 0101 1110 1000
	op := uint16(0b1001010111101000)
	return op
}
