package avr

import (
	"fmt"
)

type MCU struct {
	Name          string
	Signature     uint16
	FlashPageSize uint16
	FlashSize     uint16
	DwenBit       byte
}

var mcus = []*MCU{
	// NOTE: I'm only listing here the MCUs I own.
	//       To add new devices please open a pull request.
	{"ATtiny24", 0x910b, 0x20, 0x0800, 1 << 6},
	{"ATtiny25", 0x9108, 0x20, 0x0800, 1 << 6},
	{"ATtiny44", 0x9207, 0x40, 0x1000, 1 << 6},
	{"ATtiny45", 0x9206, 0x40, 0x1000, 1 << 6},
	{"ATtiny84", 0x930c, 0x40, 0x2000, 1 << 6},
	{"ATtiny85", 0x930b, 0x40, 0x2000, 1 << 6},
	{"ATtiny261", 0x910c, 0x20, 0x0800, 1 << 6},
	{"ATtiny861", 0x930d, 0x40, 0x2000, 1 << 6},
	{"ATtiny2313", 0x910a, 0x20, 0x0800, 1 << 7},
	{"ATtiny4313", 0x920d, 0x40, 0x1000, 1 << 7},
	{"ATmega48P", 0x920a, 0x40, 0x1000, 1 << 6},
	{"ATmega88P", 0x930f, 0x40, 0x2000, 1 << 6},
	{"ATmega328P", 0x950f, 0x80, 0x8000, 1 << 6},
}

func GetMCU(sign uint16) (*MCU, error) {
	for _, mcu := range mcus {
		if sign == mcu.Signature {
			return mcu, nil
		}
	}
	return nil, fmt.Errorf(`debugwire: failed to detect MCU from signature: 0x%04x
Please open an issue/pull request: https://github.com/dwtk/dwtk`, sign)
}
