package gdbserver

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.rgm.io/dwtk/debugwire"
)

type detachErr struct{}

func (d *detachErr) Error() string {
	return ""
}

func handleCommand(ctx context.Context, dw *debugwire.DebugWire, conn *tcpConn, cmd []byte) error {
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

	case 'G':
		cache, err := dw.Cache()
		if err != nil {
			return err
		}
		defer cache.Restore()

		b := make([]byte, hex.DecodedLen(len(cmd[1:])))
		n, err := hex.Decode(b, cmd[1:])
		if err != nil {
			return err
		}
		if n != 39 {
			return fmt.Errorf("gdbserver: commands: malformed register write request: %s", cmd)
		}

		cache.Registers = b[0:32]
		cache.SREG = b[32]
		cache.SP = uint16(b[33]) | uint16(b[34]<<8)
		cache.PC = uint16(b[35]) | uint16(b[36]<<8)

		return writePacket(conn, []byte("OK"))

	case 'g':
		cache, err := dw.Cache()
		if err != nil {
			return err
		}
		defer cache.Restore()

		b := append(
			cache.Registers,
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
		rv, err := wait(ctx, dw, conn)
		if err != nil {
			return err
		}
		return writePacket(conn, rv)

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
		if k != 2 {
			return fmt.Errorf("gdbserver: commands: invalid breakpoint size: %d", k)
		}

		switch p[0][0] {
		case '0':
			cache, err := dw.Cache()
			if err != nil {
				return err
			}
			defer cache.Restore()

			if add {
				err = dw.SetSwBreakpoint(uint16(a))
			} else {
				err = dw.ClearSwBreakpoint(uint16(a))
			}
			if err != nil {
				return err
			}

			return writePacket(conn, []byte("OK"))

		case '1':
			if add {
				if !dw.SetHwBreakpoint(uint16(a)) {
					return writePacket(conn, []byte("E01"))
				}
			} else {
				dw.ClearHwBreakpoint()
			}

			return writePacket(conn, []byte("OK"))

		default:
			return writePacket(conn, []byte("E01"))
		}

	case 'D':
		if err := writePacket(conn, []byte("OK")); err != nil {
			return err
		}
		return &detachErr{}

	case '?':
		return writePacket(conn, []byte("S00"))
	}
	return writePacket(conn, []byte{})
}
