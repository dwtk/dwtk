package usbserial

import (
	"sync"
)

type UsbSerial struct {
	Fd int

	mutex *sync.Mutex
	buf   []byte
}

func New(portDevice string, baudrate uint32) (*UsbSerial, error) {
	fd, err := Open(portDevice, baudrate)
	if err != nil {
		return nil, err
	}

	return &UsbSerial{
		Fd:    fd,
		mutex: &sync.Mutex{},
		buf:   []byte{},
	}, nil
}

func (u *UsbSerial) Commit() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	err := Write(u.Fd, u.buf)
	u.buf = []byte{}
	return err
}

func (u *UsbSerial) Close() error {
	if err := u.Commit(); err != nil {
		return err
	}

	return Close(u.Fd)
}

func (u *UsbSerial) Flush() error {
	if err := u.Commit(); err != nil {
		return err
	}

	return Flush(u.Fd)
}

func (u *UsbSerial) Read(p []byte) error {
	if err := u.Commit(); err != nil {
		return err
	}

	return Read(u.Fd, p)
}

func (u *UsbSerial) ReadByte() (byte, error) {
	var b [1]byte
	if err := u.Read(b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func (u *UsbSerial) ReadWord() (uint16, error) {
	var b [2]byte
	if err := u.Read(b[:]); err != nil {
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

func (u *UsbSerial) SendBreak() (byte, error) {
	if err := u.Commit(); err != nil {
		return 0, err
	}

	return SendBreak(u.Fd)
}

func (u *UsbSerial) RecvBreak() (byte, error) {
	if err := u.Commit(); err != nil {
		return 0, err
	}

	return RecvBreak(u.Fd)
}
