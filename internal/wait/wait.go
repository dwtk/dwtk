package wait

import (
	"context"
	"fmt"

	"golang.org/x/sys/unix"
)

func ForFd(ctx context.Context, fd int, c chan bool) error {
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

		if _, err := unix.Select(fd+1, fds, nil, nil, &unix.Timeval{Usec: 100000}); err != nil && err != unix.EINTR {
			return err
		}
		if fds.IsSet(fd) {
			c <- true
			break
		}
	}

	return nil
}
