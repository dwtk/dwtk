package debugwire

import (
	"context"
)

func (dw *DebugWire) Wait(ctx context.Context, c chan bool) error {
	return dw.Port.Wait(ctx, c)
}
