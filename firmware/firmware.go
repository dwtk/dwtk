package firmware

import (
	"fmt"
	"math"
	"os"

	"golang.rgm.io/dwtk/avr"
	"golang.rgm.io/dwtk/firmware/elf"
	"golang.rgm.io/dwtk/firmware/hex"
)

type Page struct {
	Address uint16
	Data    []byte
}

type Firmware struct {
	MCU   *avr.MCU
	Data  []byte
	Pages []*Page
}

func split(data []byte, mcu *avr.MCU) []*Page {
	pages := []*Page{}
	n := uint16(math.Ceil(float64(len(data)) / float64(mcu.FlashPageSize)))
	for i := uint16(0); i < n; i++ {
		c := make([]byte, mcu.FlashPageSize)
		addr := i * mcu.FlashPageSize
		for j := uint16(0); j < mcu.FlashPageSize && (addr+j) < uint16(len(data)); j++ {
			c[j] = data[addr+j]
		}
		pages = append(pages, &Page{
			Address: addr,
			Data:    c,
		})
	}
	return pages
}

func Empty(mcu *avr.MCU) (*Firmware, error) {
	if mcu == nil {
		return nil, fmt.Errorf("firmware: MCU must be set")
	}

	data := make([]byte, mcu.FlashSize)

	return &Firmware{
		MCU:   mcu,
		Data:  data,
		Pages: split(data, mcu),
	}, nil
}

func Parse(path string, mcu *avr.MCU) (*Firmware, error) {
	if mcu == nil {
		return nil, fmt.Errorf("firmware: MCU must be set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	data, err := func() ([]byte, error) {
		if elf.Check(path) {
			data, err := elf.Parse(path)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		if hex.Check(path) {
			data, err := hex.Parse(path)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		return nil, fmt.Errorf("firmware: failed to detect firmware file format: %s", path)
	}()
	if err != nil {
		return nil, err
	}

	if uint16(len(data)) > mcu.FlashSize {
		return nil, fmt.Errorf("firmware: size (%d) bigger than %s flash (%d)",
			len(data),
			mcu.Name,
			mcu.FlashSize,
		)
	}

	return &Firmware{
		MCU:   mcu,
		Data:  data,
		Pages: split(data, mcu),
	}, nil
}

func Dump(path string, data []byte) error {
	return hex.Dump(path, data)
}
