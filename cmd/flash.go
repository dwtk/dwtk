package cmd

import (
	"bytes"
	"flag"
	"fmt"

	"golang.rgm.io/dwtk/firmware"
	"golang.rgm.io/dwtk/internal/cli"
)

var (
	noVerify bool
)

var FlashCmd = &cli.Command{
	Name:        "flash",
	Usage:       "FILE",
	Description: "flash firmware (ELF or Intel HEX) to target MCU, verify and exit",
	SetFlags: func(fs *flag.FlagSet) error {
		fs.BoolVar(
			&noVerify,
			"n",
			false,
			"do not verify flashed firmware",
		)
		return nil
	},
	Run: func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("FILE argument required")
		}

		f, err := firmware.Parse(args[0])
		if err != nil {
			return err
		}

		pages, err := dw.MCU.PrepareFirmware(f)
		if err != nil {
			return err
		}

		i := 1
		for addr, data := range pages {
			fmt.Printf("Flashing page %d/%d ...\n", i, len(pages))
			if err := dw.WriteFlashPage(addr, data); err != nil {
				return err
			}
			i += 1
		}

		if noVerify {
			return nil
		}

		i = 1
		read := make([]byte, dw.MCU.FlashPageSize)
		for addr, data := range pages {
			fmt.Printf("Verifying page %d/%d ...\n", i, len(pages))
			if err := dw.ReadFlash(addr, read); err != nil {
				return err
			}
			if bytes.Compare(data, read) != 0 {
				return fmt.Errorf("Page mismatch 0x%04x: %v != %v", addr, data, read)
			}
			i += 1
		}

		return nil
	},
}
