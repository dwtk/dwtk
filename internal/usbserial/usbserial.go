package usbserial

import (
	"context"
	"sync"

	"github.com/dwtk/dwtk/wait"
)

type UsbSerial struct {
	baudrate uint32
	fd       int
	mutex    *sync.RWMutex
	buf      []byte
}

func Open(device string, baudrate uint32) (*UsbSerial, error) {
	fd, err := open(device, baudrate)
	if err != nil {
		return nil, err
	}

	return &UsbSerial{
		baudrate: baudrate,
		fd:       fd,
		mutex:    &sync.RWMutex{},
		buf:      []byte{},
	}, nil
}

func (u *UsbSerial) Commit() error {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	err := write(u.fd, u.buf)
	u.buf = []byte{}
	return err
}

func (u *UsbSerial) Close() error {
	if err := u.Commit(); err != nil {
		return err
	}

	return _close(u.fd)
}

func (u *UsbSerial) Flush() error {
	if err := u.Commit(); err != nil {
		return err
	}

	return flush(u.fd)
}

func (u *UsbSerial) Read(p []byte) error {
	if err := u.Commit(); err != nil {
		return err
	}

	return read(u.fd, p)
}

func (u *UsbSerial) ReadByte() (byte, error) {
	b := make([]byte, 1)
	if err := u.Read(b); err != nil {
		return 0, err
	}
	return b[0], nil
}

func (u *UsbSerial) ReadWord() (uint16, error) {
	b := make([]byte, 2)
	if err := u.Read(b); err != nil {
		return 0, err
	}
	return (uint16(b[0]) << 8) | uint16(b[1]), nil
}

func (u *UsbSerial) Write(p []byte) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.buf = append(u.buf, p...)
	return nil
}

func (u *UsbSerial) SendBreak() error {
	if err := u.Commit(); err != nil {
		return err
	}

	return sendBreak(u.fd, u.baudrate)
}

func (u *UsbSerial) RecvBreak() (byte, error) {
	if err := u.Commit(); err != nil {
		return 0, err
	}

	return recvBreak(u.fd)
}

func (u *UsbSerial) Wait(ctx context.Context, c chan bool) error {
	if err := u.Commit(); err != nil {
		return err
	}

	u.mutex.Lock()
	defer u.mutex.Unlock()

	return wait.ForFd(ctx, u.fd, c)
}
