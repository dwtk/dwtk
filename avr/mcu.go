package avr

import (
	"fmt"
	"math"
)

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
	{"ATtiny24", 0x910b, 0x20, 0x0800, 0x37},
	{"ATtiny25", 0x9108, 0x20, 0x0800, 0x37},
	{"ATtiny84", 0x930c, 0x40, 0x2000, 0x37},
	{"ATtiny85", 0x930b, 0x40, 0x2000, 0x37},
	{"ATtiny261", 0x910c, 0x20, 0x0800, 0x37},
	{"ATtiny861", 0x930d, 0x40, 0x2000, 0x37},
	{"ATtiny2313", 0x910a, 0x20, 0x0800, 0x37},
	{"ATtiny4313", 0x920d, 0x40, 0x1000, 0x37},
	{"ATmega48A", 0x920a, 0x40, 0x1000, 0x37},
	{"ATmega88A", 0x930f, 0x40, 0x2000, 0x37},
	{"ATmega328", 0x950f, 0x80, 0x8000, 0x37},
}

func GetMCU(sign uint16) *MCU {
	for _, mcu := range mcus {
		if sign == mcu.Signature {
			return mcu
		}
	}
	return nil
}

func (m *MCU) PrepareFirmware(data []byte) (map[uint16][]byte, error) {
	if uint16(len(data)) > m.FlashSize {
		return nil, fmt.Errorf("mcu: firmware size (%d) bigger than %s flash (%d)",
			len(data),
			m.Name,
			m.FlashSize,
		)
	}

	rv := make(map[uint16][]byte)
	n := uint16(math.Ceil(float64(len(data)) / float64(m.FlashPageSize)))
	for i := uint16(0); i < n; i++ {
		c := make([]byte, m.FlashPageSize)
		addr := i * m.FlashPageSize
		for j := uint16(0); j < m.FlashPageSize && (addr+j) < uint16(len(data)); j++ {
			c[j] = data[addr+j]
		}
		rv[addr] = c
	}

	return rv, nil
}

func (m *MCU) NumFlashPages() (uint16, error) {
	if m.FlashSize%m.FlashPageSize != 0 {
		return 0, fmt.Errorf("avr: invalid flash size: 0x%04x", m.FlashSize)
	}

	return m.FlashSize / m.FlashPageSize, nil
}
