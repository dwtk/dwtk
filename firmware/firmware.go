package firmware

import (
	"fmt"
	"os"

	"golang.rgm.io/dwtk/firmware/elf"
	"golang.rgm.io/dwtk/firmware/hex"
)

func Parse(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

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
}
