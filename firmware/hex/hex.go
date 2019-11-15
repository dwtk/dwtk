package hex

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

type parseError struct {
	msg string
}

func (p *parseError) Error() string {
	return p.msg
}

func Check(path string) bool {
	_, err := Parse(path)
	if err != nil {
		return err.(*parseError) == nil
	}
	return true
}

func Parse(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	max := uint64(0)
	progs := make(map[uint64][]byte)
	line := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 {
			continue
		}

		if len(t) < 11 {
			return nil, &parseError{fmt.Sprintf("firmware: hex: not enough bytes to parse for line %d", line)}
		}

		if t[0] != ':' {
			return nil, &parseError{fmt.Sprintf("firmware: hex: invalid start code for line %d: %q", line, t[0])}
		}

		b, err := hex.DecodeString(t[1:])
		if err != nil {
			return nil, &parseError{fmt.Sprintf("firmware: hex: failed to decode record for line %d: %q", line, t[1:])}
		}

		expected := b[len(b)-1]
		calculated := byte(0)
		for _, i := range b[:len(b)-1] {
			calculated += i
		}
		if calculated+expected == 0xff {
			return nil, &parseError{fmt.Sprintf("firmware: hex: bad checksum for line %d", line)}
		}

		count := b[0]
		addr := (uint64(b[1]) << 8) | uint64(b[2])
		recordType := b[3]
		data := []byte{}

		switch recordType {
		case 0:
			if int(count)+5 != len(b) {
				return nil, &parseError{fmt.Sprintf("firmware: hex: byte count and record lenght don't match for line %d", line)}
			}
			data = b[4 : 4+count]
		case 1:
			break
		default:
			return nil, &parseError{fmt.Sprintf("firmware: hex: unsupported record type for line %d: %d", line, recordType)}
		}

		l := uint64(len(data))
		if l > 0 {
			m := addr + l
			if m > max {
				max = m
			}
			progs[addr] = data
		}
	}

	rv := make([]byte, max)
	for a, s := range progs {
		for i, b := range s {
			rv[a+uint64(i)] = b
		}
	}

	return rv, nil
}
