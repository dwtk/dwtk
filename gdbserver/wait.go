package gdbserver

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
	"golang.rgm.io/dwtk/debugwire"
)

func wait(dw *debugwire.DebugWire, conn net.Conn) (bool, error) {
	fds := &unix.FdSet{}
	fds.Zero()

	nfds := 0
	if dw.Port.Fd >= 0 {
		fds.Set(dw.Port.Fd)
		nfds = dw.Port.Fd
	}

	c := conn.(*net.TCPConn)
	if c == nil {
		return false, fmt.Errorf("gdbserver: wait: invalid tcp connection")
	}
	f, err := c.File()
	if err != nil {
		return false, err
	}
	fd := int(f.Fd())
	if fd >= 0 {
		fds.Set(fd)
		if fd > nfds {
			nfds = fd
		}
	}

	r, err := unix.Select(nfds+1, fds, nil, nil, nil)
	if err != nil {
		return false, err
	}
	if r == -1 {
		return false, fmt.Errorf("gdbserver: wait: failed select")
	}
	if r == 0 {
		return false, fmt.Errorf("gdbserver: wait: failed select, no data")
	}
	if fds.IsSet(dw.Port.Fd) {
		return false, dw.RecvBreak()
	}

	return true, nil
}
