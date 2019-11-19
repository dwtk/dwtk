package gdbserver

import (
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
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, " * Connection accepted from %s\n", conn.RemoteAddr().String())

	if err := dw.Reset(); err != nil {
		return err
	}

	defer func() {
		if err := dw.ClearSwBreakpoints(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to clear software breakpoints: %s\n", err)
		}
	}()

	for {
		if err := handlePacket(dw, conn); err != nil {
			if _, ok := err.(*detachErr); ok {
				return nil
			}
			return err
		}
	}

	return nil
}
