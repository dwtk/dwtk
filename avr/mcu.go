package avr

type MCU struct {
	Name          string
	Signature     uint16
	FlashPageSize uint16
	FlashSize     uint16
	SPMCSR        byte
}

var mcus = []*MCU{
	// NOTE: I'm only listing here the MCUs I own.
	//       To add new devices please open a pull request.
	{"ATtiny24/A", 0x910b, 0x20, 0x0800, 0x37},
	{"ATtiny25", 0x9108, 0x20, 0x0800, 0x37},
	{"ATtiny84/A", 0x930c, 0x40, 0x2000, 0x37},
	{"ATtiny85", 0x930b, 0x40, 0x2000, 0x37},
	{"ATtiny261/A", 0x910c, 0x20, 0x0800, 0x37},
	{"ATtiny861/A", 0x930d, 0x40, 0x2000, 0x37},
	{"ATtiny2313/A", 0x910a, 0x20, 0x0800, 0x37},
	{"ATtiny4313", 0x920d, 0x40, 0x1000, 0x37},
	{"ATmega48A/P/PA", 0x920a, 0x40, 0x1000, 0x37},
	{"ATmega88A/P/PA", 0x930f, 0x40, 0x2000, 0x37},
	{"ATmega328/P", 0x950f, 0x80, 0x8000, 0x37},
}

func GetMCU(sign uint16) *MCU {
	for _, mcu := range mcus {
		if sign == mcu.Signature {
			return mcu
		}
	}
	return nil
}
