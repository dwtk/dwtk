package hex

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
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
	line := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line++

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
		if calculated+expected != 0x00 {
			return nil, &parseError{fmt.Sprintf("firmware: hex: bad checksum for line %d", line)}
		}

		count := b[0]
		addr := (uint64(b[1]) << 8) | uint64(b[2])
		recordType := b[3]
		data := []byte{}

		switch recordType {
		case 0:
			if int(count)+5 != len(b) {
				return nil, &parseError{fmt.Sprintf("firmware: hex: byte count and record length don't match for line %d", line)}
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

func Dump(path string, data []byte) error {
	rv := []byte{}

	for i := uint16(0); i < uint16(len(data)); i += 16 {
		end := i + 16
		if end > uint16(len(data)) {
			end = uint16(len(data))
		}

		d := data[i:end]
		b := append([]byte{byte(len(d)), byte(i >> 8), byte(i), 0}, d...)

		// not using hex.Encode because we want upper case
		chk := byte(0)
		rv = append(rv, ':')
		for _, j := range b {
			chk += j
			rv = append(rv, []byte(fmt.Sprintf("%02X", j))...)
		}

		rv = append(rv, []byte(fmt.Sprintf("%02X\n", byte(256-uint16(chk))))...)
	}

	rv = append(rv, ':', '0', '0', '0', '0', '0', '0', '0', '1', 'F', 'F')

	return ioutil.WriteFile(path, rv, 0666)
}
