package hex

import (
	"path"

	"github.com/dwtk/dwtk/internal/hex"
)

type Hex struct{}

func (*Hex) Check(fpath string) bool {
	if path.Ext(fpath) == ".hex" {
		return true
	}
	_, err := hex.Parse(fpath)
	return err == nil
}

func (*Hex) Parse(fpath string) ([]byte, error) {
	return hex.Parse(fpath)
}

func (*Hex) Dump(fpath string, data []byte) error {
	return hex.Dump(fpath, data)
}
