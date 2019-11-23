package gdbserver

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func signalChannel() chan os.Signal {
	sig := make(chan os.Signal)
	signal.Notify(sig, unix.SIGINT, unix.SIGKILL, unix.SIGTERM)
	return sig
}

func waitForFd(ctx context.Context, fd int, c chan bool) error {
	if fd < 0 {
		return fmt.Errorf("bad file descriptor: %d", fd)
	}

	fds := &unix.FdSet{}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fds.Zero()
		fds.Set(fd)

		r, err := unix.Select(fd+1, fds, nil, nil, &unix.Timeval{Usec: 100000})
		if err != nil && err != unix.EINTR {
			return err
		}
		if r == -1 {
			return fmt.Errorf("failed select")
		}
		if fds.IsSet(fd) {
			c <- true
			break
		}
	}

	return nil
}
