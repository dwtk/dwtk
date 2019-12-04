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
	Data []byte
	MCU  *avr.MCU
}

func NewEmpty(mcu *avr.MCU) (*Firmware, error) {
	if mcu == nil {
		return nil, fmt.Errorf("firmware: MCU must be set")
	}

	data := make([]byte, mcu.FlashSize)
	for i := range data {
		data[i] = 0xff
	}

	// not using NewFromData because we know that our "firmware" size is safe
	return &Firmware{
		MCU:  mcu,
		Data: data,
	}, nil
}

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

	return NewFromData(data, mcu)
}

func (f *Firmware) SplitPages() []*Page {
	pages := []*Page{}
	n := uint16(math.Ceil(float64(len(f.Data)) / float64(f.MCU.FlashPageSize)))
	for i := uint16(0); i < n; i++ {
		c := make([]byte, f.MCU.FlashPageSize)
		addr := i * f.MCU.FlashPageSize
		for j := uint16(0); j < f.MCU.FlashPageSize && (addr+j) < uint16(len(f.Data)); j++ {
			c[j] = f.Data[addr+j]
		}
		pages = append(pages, &Page{
			Address: addr,
			Data:    c,
		})
	}
	return pages
}

func (f *Firmware) Dump(path string) error {
	return hex.Dump(path, f.Data)
}
