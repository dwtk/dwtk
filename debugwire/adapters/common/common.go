package common

import (
	"github.com/dwtk/devices"
)

type Common interface {
	GetMCU() *devices.MCU
	WriteRegisters(start byte, regs []byte) error
	ReadRegisters(start byte, regs []byte) error
	WriteInstruction(inst uint16) error
}
