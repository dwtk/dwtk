package elf

import (
	"debug/elf"
	"fmt"
	"io/ioutil"
	"path"
)

type ELF struct{}

func (*ELF) Check(fpath string) bool {
	if path.Ext(fpath) == ".elf" {
		return true
	}
	f, err := elf.Open(fpath)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func (*ELF) Parse(fpath string) ([]byte, error) {
	f, err := elf.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if f.Machine != elf.EM_AVR {
		return nil, fmt.Errorf("firmware: elf: invalid machine architecture: %s", f.Machine)
	}

	max := uint64(0)
	progs := make(map[uint64][]byte)
	for _, s := range f.Progs {
		if s.Type != elf.PT_LOAD {
			continue
		}

		data, err := ioutil.ReadAll(s.Open())
		if err != nil {
			return nil, err
		}

		addr := uint64(s.Paddr)
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
