package firmware

import (
	"fmt"
	"math"
	"os"

	"github.com/dwtk/dwtk/avr"
	"github.com/dwtk/dwtk/firmware/elf"
	"github.com/dwtk/dwtk/firmware/hex"
)

type Page struct {
	Address uint16
	Data    []byte
}

type Firmware struct {
	Data []byte
	MCU  *avr.MCU
}

type format interface {
	Check(fpath string) bool
	Parse(fpath string) ([]byte, error)
}

var (
	formats = []format{
		&elf.ELF{},
		&hex.Hex{},
	}
)

func NewFromData(data []byte, mcu *avr.MCU) (*Firmware, error) {
	if mcu == nil {
		return nil, fmt.Errorf("firmware: MCU must be set")
	}

	if uint16(len(data)) > mcu.FlashSize {
		return nil, fmt.Errorf("firmware: size (%d) bigger than %s flash (%d)",
			len(data),
			mcu.Name,
			mcu.FlashSize,
		)
	}

	return &Firmware{
		Data: data,
		MCU:  mcu,
	}, nil
}

func NewFromFile(path string, mcu *avr.MCU) (*Firmware, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	data, err := func() ([]byte, error) {
		for _, format := range formats {
			if format.Check(path) {
				data, err := format.Parse(path)
				if err != nil {
					return nil, err
				}
				return data, nil
			}
		}
		return nil, fmt.Errorf("firmware: failed to detect firmware file format: %s", path)
	}()
	if err != nil {
		return nil, err
	}

	return NewFromData(data, mcu)
}

func (f *Firmware) SplitPages() []*Page {
	pages := []*Page{}
	n := uint16(math.Ceil(float64(len(f.Data)) / float64(f.MCU.FlashPageSize)))
	for i := uint16(0); i < n; i++ {
		c := make([]byte, f.MCU.FlashPageSize)
		addr := i * f.MCU.FlashPageSize
		var j uint16
		for j = 0; j < f.MCU.FlashPageSize && (addr+j) < uint16(len(f.Data)); j++ {
			c[j] = f.Data[addr+j]
		}
		for ; j < f.MCU.FlashPageSize; j++ {
			c[j] = 0xff // empty = 0xff
		}
		pages = append(pages, &Page{
			Address: addr,
			Data:    c,
		})
	}
	return pages
}

func (f *Firmware) Dump(path string) error {
	h := &hex.Hex{}
	return h.Dump(path, f.Data)
}
