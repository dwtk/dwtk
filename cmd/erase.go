package cmd

import (
	"fmt"

	"golang.rgm.io/dwtk/internal/cli"
)

var EraseCmd = &cli.Command{
	Name:        "erase",
	Description: "erase the target MCU's flash and exit",
	Run: func(args []string) error {
		noReset = true

		pages, err := dw.MCU.PrepareFirmware(make([]byte, dw.MCU.FlashSize))
		if err != nil {
			return err
		}

		i := 1
		for addr, data := range pages {
			fmt.Printf("Erasing page %d/%d ...\n", i, len(pages))
			if err := dw.WriteFlashPage(addr, data); err != nil {
				return err
			}
			i += 1
		}

		return nil
	},
}
