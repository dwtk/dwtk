package gdbserver

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dwtk/dwtk/debugwire"
)

type detachErr struct{}

func (d *detachErr) Error() string {
	return ""
}

func handleCommand(ctx context.Context, dw *debugwire.DebugWIRE, conn *tcpConn, cmd []byte) error {
	notifyGdb := func(err error, rsp []byte) error {
		errs := []string{}
		if err != nil {
			errs = append(errs, err.Error())
			if e := writePacket(conn, rsp); e != nil {
				errs = append(errs, e.Error())
			}
		}

		if len(errs) > 0 {
			return errors.New(strings.Join(errs, "; "))
		}

		return nil
	}

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
		b := make([]byte, hex.DecodedLen(len(cmd[1:])))
		n, err := hex.Decode(b, cmd[1:])
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		if n != 39 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed register write request: %s", cmd),
				[]byte("E01"),
			)
		}

		if err := dw.WriteRegisters(0, b[0:32]); err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		if err := dw.SetSREG(b[32]); err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		if err := dw.SetSP(uint16(b[33]) | (uint16(b[34]) << 8)); err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		if err := dw.SetPC(uint16(b[35]) | (uint16(b[36]) << 8)); err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		return writePacket(conn, []byte("OK"))

	case 'g':
		b := make([]byte, 32)
		if err := dw.ReadRegisters(0, b); err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		sreg, err := dw.GetSREG()
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		sp, err := dw.GetSP()
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		pc, err := dw.GetPC()
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		b = append(b,
			sreg,
			byte(sp), byte(sp>>8),
			byte(pc), byte(pc>>8), 0, 0,
		)

		d := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(d, b)
		return writePacket(conn, d)

	case 'P':
		p := strings.Split(scmd[1:], "=")
		if len(p) != 2 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed register write request: %s", cmd),
				[]byte("E01"),
			)
		}

		a, err := strconv.ParseUint(p[0], 16, 8)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		b, err := strconv.ParseUint(p[1], 16, 32)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		switch a {
		case 32:
			if err := dw.SetSREG(byte(b)); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		case 33:
			if err := dw.SetSP(uint16(b>>8) | (uint16(b&0xff) << 8)); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		case 34:
			if err := dw.SetPC(uint16(b>>24) | (uint16(b >> 8))); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		default:
			if err := dw.WriteRegisters(byte(a), []byte{byte(b)}); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		}

		return writePacket(conn, []byte("OK"))

	case 'p':
		a, err := strconv.ParseUint(scmd[1:], 16, 8)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		b := []byte{}

		switch a {
		case 32:
			sreg, err := dw.GetSREG()
			if err != nil {
				return notifyGdb(err, []byte("E01"))
			}
			b = append(b, sreg)
		case 33:
			sp, err := dw.GetSP()
			if err != nil {
				return notifyGdb(err, []byte("E01"))
			}
			b = append(b, byte(sp), byte(sp>>8))
		case 34:
			pc, err := dw.GetPC()
			if err != nil {
				return notifyGdb(err, []byte("E01"))
			}
			b = append(b, byte(pc), byte(pc>>8), 0, 0)
		default:
			t := make([]byte, 1)
			if err := dw.ReadRegisters(byte(a), t); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
			b = append(b, t...)
		}

		d := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(d, b)
		return writePacket(conn, d)

	case 'M':
		h := strings.Split(scmd[1:], ":")
		if len(h) != 2 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed memory write request: %s", cmd),
				[]byte("E01"),
			)

		}

		p := strings.Split(h[0], ",")
		if len(p) != 2 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed memory write request: %s", cmd),
				[]byte("E01"),
			)
		}

		a, err := strconv.ParseUint(p[0], 16, 32)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		c, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		b, err := hex.DecodeString(h[1])
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		if uint64(len(b)) != c {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed memory write request: %s", cmd),
				[]byte("E01"),
			)
		}

		if a < 0x800000 {
			if err := dw.WriteFlash(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		} else if a < 0x810000 {
			if err := dw.WriteSRAM(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		} else {
			if err := dw.WriteEEPROM(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		}

		return writePacket(conn, []byte("OK"))

	case 'm':
		p := strings.Split(scmd[1:], ",")
		if len(p) != 2 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed memory read request: %s", cmd),
				[]byte("E01"),
			)
		}

		a, err := strconv.ParseUint(p[0], 16, 32)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		c, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		b := make([]byte, c)
		if a < 0x800000 {
			if err := dw.ReadFlash(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		} else if a < 0x810000 {
			if err := dw.ReadSRAM(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		} else {
			if err := dw.ReadEEPROM(uint16(a), b); err != nil {
				return notifyGdb(err, []byte("E01"))
			}
		}

		d := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(d, b)
		return writePacket(conn, d)

	case 's':
		if err := dw.Step(); err != nil {
			return notifyGdb(err, []byte("S00"))
		}
		return writePacket(conn, []byte("S05"))

	case 'c':
		if err := dw.Continue(); err != nil {
			return notifyGdb(err, []byte("S00"))
		}
		rv, err := waitForDwOrGdb(ctx, dw, conn)
		if err != nil {
			return notifyGdb(err, []byte("S00"))
		}
		return writePacket(conn, rv)

	case 'Z':
		add = true
		fallthrough

	case 'z':
		p := strings.Split(scmd[1:], ",")
		if len(p) != 3 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: malformed breakpoint request: %s", cmd),
				[]byte("E01"),
			)
		}

		a, err := strconv.ParseUint(p[1], 16, 16)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}

		k, err := strconv.ParseUint(p[2], 16, 8)
		if err != nil {
			return notifyGdb(err, []byte("E01"))
		}
		if k != 2 {
			return notifyGdb(
				fmt.Errorf("gdbserver: commands: invalid breakpoint size: %d", k),
				[]byte("E01"),
			)
		}

		switch p[0][0] {
		case '0':
			if add {
				err = dw.SetSwBreakpoint(uint16(a))
			} else {
				err = dw.ClearSwBreakpoint(uint16(a))
			}
			if err != nil {
				return notifyGdb(err, []byte("E01"))
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
			return notifyGdb(err, []byte("E01"))
		}
		return &detachErr{}

	case '?':
		return writePacket(conn, []byte("S00"))
	}

	return writePacket(conn, []byte{})
}
