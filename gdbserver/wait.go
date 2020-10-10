package gdbserver

import (
	"context"
	"fmt"
	"os"

	"github.com/dwtk/dwtk/debugwire"
	"github.com/dwtk/dwtk/internal/wait"
)

func waitForDwOrGdb(ctx context.Context, dw *debugwire.DebugWIRE, conn *tcpConn) ([]byte, error) {
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigGdb := make(chan bool)
	sigDw := make(chan bool)

	go func() {
		if err := wait.ForFd(nctx, conn.Fd, sigGdb); err != nil {
			fmt.Fprintf(os.Stderr, "error: gdbserver: gdb: %s\n", err)
		}
	}()

	go func() {
		if err := dw.Wait(nctx, sigDw); err != nil {
			fmt.Fprintf(os.Stderr, "error: gdbserver: debugwire: %s\n", err)
		}
	}()

	select {
	case <-ctx.Done():
		return []byte("S00"), nil
	case <-sigGdb:
		return []byte("S02"), nil
	case <-sigDw:
		return []byte("S05"), dw.RecvBreak()
	}
}
