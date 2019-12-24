package usbserial

import (
	"context"
)

func (us *UsbSerialAdapter) Wait(ctx context.Context, c chan bool) error {
	return us.device.Wait(ctx, c)
}
