package dwtk

import (
	"context"
	"time"
)

func (dw *DwtkAdapter) Wait(ctx context.Context, c chan bool) error {
	f := make([]byte, 1)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := dw.controlIn(cmdWait, 0, 0, f); err != nil {
			return err
		}

		if f[0] != 0 {
			c <- true
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
