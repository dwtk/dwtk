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

		n, err := unix.Select(fd+1, fds, nil, nil, &unix.Timeval{Usec: 100000})
		if err != nil && err != unix.EINTR {
			return err
		}
		if n == 1 && fds.IsSet(fd) {
			c <- true
			break
		}
	}

	return nil
}
