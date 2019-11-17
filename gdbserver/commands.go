package gdbserver

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"golang.rgm.io/dwtk/avr"
	"golang.rgm.io/dwtk/debugwire"
)

func handleCommand(dw *debugwire.DebugWire, conn net.Conn, cmd []byte) error {
	if len(cmd) == 0 {
		return fmt.Errorf("gdbserver: commands: empty command")
	}
	scmd := string(cmd)

	add := false
	switch cmd[0] {
	case 'q':
		if scmd == "qAttached" {
			return writePacket(conn, []byte{'1'})
		}

	case 'g':
		cache, err := dw.Cache()
		if err != nil {
			return err
		}
		defer cache.Restore()

		b := append(
			cache.Registers[:],
			cache.SREG,
			byte(cache.SP), byte(cache.SP>>8),
			byte(cache.PC), byte((cache.PC)>>8), 0, 0,
		)
		d := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(d, b)
		return writePacket(conn, d)

	case 'm':
		cache, err := dw.Cache()
		if err != nil {
			return err
		}
		defer cache.Restore()

		p := strings.Split(scmd[1:], ",")
		if len(p) != 2 {
			writePacket(conn, []byte("E01"))
			return fmt.Errorf("gdbserver: commands: malformed memory read request: %s", cmd)
		}

		a, err := strconv.ParseUint(p[0], 16, 32)
		if err != nil {
			return err
		}

		c, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			return err
		}

		b := make([]byte, c)
		if a < 0x800000 {
			if err := dw.ReadFlash(uint16(a), b); err != nil {
				return err
			}
		} else if a < 0x810000 {
			if err := dw.ReadSRAM(uint16(a), b); err != nil {
				return err
			}
		} else { // eeprom
			writePacket(conn, []byte("E01"))
			return nil
		}

		d := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(d, b)
		return writePacket(conn, d)

	case 's':
		if err := dw.Step(); err != nil {
			return err
		}
		return writePacket(conn, []byte("S05"))

	case 'c':
		if err := dw.Continue(); err != nil {
			return err
		}
		interrupt, err := wait(dw, conn)
		if err != nil {
			return err
		}
		if interrupt {
			return writePacket(conn, []byte("S02"))
		}
		return writePacket(conn, []byte("S05"))

	case 'Z':
		add = true
		fallthrough

	case 'z':
		p := strings.Split(scmd[1:], ",")
		if len(p) != 3 {
			writePacket(conn, []byte("E01"))
			return fmt.Errorf("gdbserver: commands: malformed breakpoint request: %s", cmd)
		}

		a, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			return err
		}
		k, err := strconv.ParseUint(p[2], 16, 8)
		if err != nil {
			return err
		}

		switch p[0][0] {
		case '0':
			cache, err := dw.Cache()
			if err != nil {
				return err
			}
			defer cache.Restore()

			addr := uint16(a)
			if add {
				f := make([]byte, 2)
				if err := dw.ReadFlash(addr, f); err != nil {
					return err
				}
				dw.SwBreakpoints[addr] = (uint16(f[1]) << 8) | uint16(f[0])
				if err := dw.WriteFlashWord(addr, avr.BREAK()); err != nil {
					// FIXME: try to recover other breakpoints
					return err
				}
			} else {
				bp, ok := dw.SwBreakpoints[addr]
				if ok {
					if err := dw.WriteFlashWord(addr, bp); err != nil {
						return err
					}
					delete(dw.SwBreakpoints, addr)
				}
			}
			return writePacket(conn, []byte("OK"))

		case '1':
			if add {
				if dw.HwBreakpointSet {
					return writePacket(conn, []byte("E01"))
				}
				dw.HwBreakpointSet = true
				dw.HwBreakpoint = uint16(a) / uint16(k)
			} else {
				dw.HwBreakpointSet = false
				dw.HwBreakpoint = 0
			}
			return writePacket(conn, []byte("OK"))

		default:
			return writePacket(conn, []byte("E01"))
		}

	case '?':
		return writePacket(conn, []byte("S00"))
	}
	return writePacket(conn, []byte{})
}
