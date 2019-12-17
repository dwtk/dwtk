package usbserial

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/sys/unix"
	"golang.rgm.io/dwtk/logger"
)

func open(portDevice string, baudrate uint32) (int, error) {
	fd, err := unix.Open(portDevice, unix.O_RDWR, 0600)
	if err != nil {
		return -1, err
	}

	cfg := &unix.Termios{
		Iflag:  unix.IGNPAR,
		Cflag:  unix.BOTHER | unix.CS8 | unix.CLOCAL,
		Oflag:  0,
		Lflag:  0,
		Ispeed: baudrate,
		Ospeed: baudrate,
	}
	cfg.Cc[unix.VMIN] = 0

	// as we always receive the echo from our own transmission, 200ms is enough, and bigger than our maximum frame time
	cfg.Cc[unix.VTIME] = 2

	if err := unix.IoctlSetTermios(fd, unix.TCSETS2, cfg); err != nil {
		unix.Close(fd)
		return -1, err
	}

	time.Sleep(30 * time.Millisecond)

	if err := flush(fd); err != nil {
		unix.Close(fd)
		return -1, err
	}

	return fd, nil
}

func _close(fd int) error {
	return unix.Close(fd)
}

func flush(fd int) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIOFLUSH)
}

func read(fd int, p []byte) error {
	n := 0
	for n < len(p) {
		c, err := unix.Read(fd, p[n:])
		if err != nil {
			return err
		}
		if c == 0 {
			return fmt.Errorf("usbserial: read: got unexpected EOF")
		}
		n += c
	}
	for i := 0; i < n; i++ {
		logger.Debug.Printf("<<< 0x%02x", p[i])
	}
	return nil
}

func write(fd int, p []byte) error {
	n := 0
	for n < len(p) {
		c, err := unix.Write(fd, p[n:])
		if err != nil {
			return err
		}
		if c == 0 {
			return fmt.Errorf("usbserial: write: got unexpected EOF")
		}
		n += c
	}
	for i := 0; i < n; i++ {
		logger.Debug.Printf(">>> 0x%02x", p[i])
	}

	e := make([]byte, len(p))
	if err := read(fd, e); err != nil {
		return err
	}

	if bytes.Compare(p, e) != 0 {
		return fmt.Errorf("usbserial: got unexpected byte echoed back")
	}

	return nil
}

func sendBreak(fd int, baudrate uint32) error {
	logger.Debug.Print("> break")

	if err := unix.IoctlSetInt(fd, unix.TIOCSBRK, 0); err != nil {
		return err
	}

	// we consider a frame as 10 bits, and send ~2 frames.
	time.Sleep(time.Duration(float64(20000000)/float64(baudrate)) * time.Microsecond)

	return unix.IoctlSetInt(fd, unix.TIOCCBRK, 0)
}

func recvBreak(fd int) (byte, error) {
	logger.Debug.Print("< break")

	c := make([]byte, 1)
	for {
		if err := read(fd, c); err != nil {
			return 0, err
		}
		if c[0] != 0x00 && c[0] != 0xff {
			break
		}
	}

	return c[0], nil
}
