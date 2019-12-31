package gdbserver

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/dwtk/dwtk/internal/wait"
)

type tcpConn struct {
	*net.TCPConn
	Fd int
}

func newConn(conn net.Conn) (*tcpConn, error) {
	c, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("gdbserver: net: invalid tcp connection")
	}

	f, err := c.File()
	if err != nil {
		return nil, err
	}
	return &tcpConn{TCPConn: c, Fd: int(f.Fd())}, nil
}

func (conn *tcpConn) readByte(ctx context.Context) (byte, error) {
	c := make(chan bool)
	go func() {
		if err := wait.ForFd(ctx, conn.Fd, c); err != nil {
			fmt.Fprintf(os.Stderr, "error: gdbserver: gdb: %s\n", err)
		}
	}()

	select {
	case <-ctx.Done():
		return 0, nil
	case <-c:
	}

	d := make([]byte, 1)
	n, err := conn.Read(d)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("gdbserver: net: failed to read byte")
	}

	return d[0], nil
}
