package gdbserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/sys/unix"
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

	sigInt := make(chan os.Signal)
	signal.Notify(sigInt, unix.SIGINT, unix.SIGKILL, unix.SIGTERM)
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
		errs := []string{}
		if errg != nil {
			errs = append(errs, errg.Error())
		}

		if err := dw.Reset(); err != nil {
			errs = append(errs, err.Error())
		} else {
			if err := dw.ClearSwBreakpoints(); err != nil {
				errs = append(errs, fmt.Sprintf("gdbserver: failed to clear software breakpoints: %s", err))
			}
		}

		if len(errs) > 0 {
			return errors.New(strings.Join(errs, "\n"))
		}
	}

	return errg
}
