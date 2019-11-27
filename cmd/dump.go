package cmd

import (
	"fmt"

	"golang.rgm.io/dwtk/firmware"
	"golang.rgm.io/dwtk/internal/cli"
)

var DumpCmd = &cli.Command{
	Name:        "dump",
	Usage:       "FILE",
	Description: "dump firmware (Intel HEX) from target MCU and exit",
	Run: func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("FILE argument required")
		}

		numPages, err := dw.MCU.NumFlashPages()
		if err != nil {
			return err
		}

		read := make([]byte, dw.MCU.FlashPageSize)
		f := []byte{}
		for i := uint16(0); i < numPages; i += 1 {
			fmt.Printf("Retrieving page %d/%d ...\n", i+1, numPages)
			addr := i * dw.MCU.FlashPageSize
			if err := dw.ReadFlash(addr, read); err != nil {
				return err
			}
			f = append(f, read...)
		}

		return firmware.Dump(args[0], f)
	},
}
