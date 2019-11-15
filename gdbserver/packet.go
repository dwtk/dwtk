package gdbserver

import (
	"encoding/hex"
	"fmt"
	"net"

	"golang.rgm.io/dwtk/debugwire"
	"golang.rgm.io/dwtk/logger"
)

type packetState uint8

const (
	packetAck packetState = iota
	packetStart
	packetCommand
	packetChecksum1
	packetChecksum2
)

func handlePacket(dw *debugwire.DebugWire, conn net.Conn) error {
	var (
		cmd  []byte
		cmdl []byte
		chk  byte
		chkr [2]byte
	)

	state := packetAck
	for {
		d := make([]byte, 1)
		n, err := conn.Read(d)
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("gdbserver: packet: failed to read byte from client socket")
		}
		b := d[0]

		if b == 0x03 {
			logger.Debug.Println("$< ctrl-c")
			if err := dw.SendBreak(); err != nil {
				return err
			}
			state = packetAck
			continue
		}

		switch state {
		case packetAck:
			if b == '+' {
				logger.Debug.Println("$< ack")
				break
			}
			if b == '-' {
				logger.Debug.Println("$< nack")
				if cmdl != nil {
					if err := handleCommand(dw, conn, cmdl); err != nil {
						return err
					}
				}
				return nil
			}
			if b == '$' {
				state = packetStart
				break
			}
			return fmt.Errorf("gdbserver: packet: ack failed, expected '+', got '%c'", b)

		case packetStart:
			cmd = []byte{b}
			chk = b
			state = packetCommand

		case packetCommand:
			if b == '#' {
				state = packetChecksum1
				break
			}
			cmd = append(cmd, b)
			chk += b

		case packetChecksum1:
			chkr[0] = b
			state = packetChecksum2

		case packetChecksum2:
			chkr[1] = b
			state = packetAck

			chkg := make([]byte, 1)
			if _, err := hex.Decode(chkg, chkr[:]); err != nil {
				return err
			}

			if chk != chkg[0] {
				return fmt.Errorf("gdbserver: packet: bad checksum, expected '0x%02x', got '0x%02x'", chkg[0], chk)
			}

			logger.Debug.Printf("$< command: %s\n", cmd)

			logger.Debug.Println("$> ack")
			n, err := conn.Write([]byte{'+'})
			if err != nil {
				return err
			}
			if n != 1 {
				return fmt.Errorf("gdbserver: packet: failed to write ack byte to client socket")
			}

			cmdl = cmd
			if err := handleCommand(dw, conn, cmd); err != nil {
				return err
			}
		}
	}

	return nil
}

func writePacket(conn net.Conn, b []byte) error {
	chk := byte(0)
	for i := 0; i < len(b); i++ {
		chk += b[i]
	}

	logger.Debug.Printf("$> command: %s\n", b)
	r := []byte(fmt.Sprintf("$%s#%02x", b, chk))

	n := 0
	for n < len(r) {
		c, err := conn.Write(r[n:])
		if err != nil {
			return err
		}
		if c == 0 {
			return fmt.Errorf("gdbserver: packet: got unexpected EOF")
		}
		n += c
	}

	return nil
}