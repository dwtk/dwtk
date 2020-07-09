package usbserial

import (
	"github.com/dwtk/dwtk/debugwire/adapters/common"
)

func (us *UsbSerialAdapter) ReadFuses() ([]byte, error) {
	return common.ReadFuses(us)
}
