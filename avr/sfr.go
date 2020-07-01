package avr

type SFR byte

func (s SFR) IO() byte {
	return byte(s)
}

func (s SFR) Mem() uint16 {
	// this is only used with WriteSRAM APIs, that take the address as a word
	return uint16(s) + 0x20
}
