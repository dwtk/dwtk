package gdbserver

import (
	"context"
	"fmt"
	"net"
	"os"

	"golang.rgm.io/dwtk/debugwire"
)

func ListenAndServe(addr string, dw *debugwire.DebugWire) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Fprintf(os.Stderr, " * GDB server running on %s\n", addr)

	// an usual socket server would loop here, accept multiple connections
	// and handle them in goroutines, but we just want to handle the first
	// incoming request.
	c, err := ln.Accept()
	if err != nil {
		return err
	}
	defer c.Close()

	conn, err := newConn(c)
	if err != nil {
		return err
	}

	sigInt := signalChannel()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigInt
		cancel()
	}()

	fmt.Fprintf(os.Stderr, " * Connection accepted from %s\n", conn.RemoteAddr().String())

	if err := dw.Reset(); err != nil {
		return err
	}

	var errg error

	for {
		exit := false
		select {
		case <-ctx.Done():
			exit = true
		default:
		}
		if exit {
			break
		}

		if err := handlePacket(ctx, dw, conn); err != nil {
			if _, ok := err.(*detachErr); ok {
				break
			}
			errg = err
			break
		}
	}

	if dw.HasSwBreakpoints() {
		if err := dw.Reset(); err != nil {
			if errg == nil {
				return err
			}
			return fmt.Errorf("%s\n%s", errg, err)
		}

		if err := dw.ClearSwBreakpoints(); err != nil {
			if errg == nil {
				return fmt.Errorf("gdbserver: failed to clear software breakpoints: %s", err)
			}
			return fmt.Errorf("%s\ngdbserver: failed to clear software breakpoints: %s", errg, err)
		}
	}

	return errg
}
