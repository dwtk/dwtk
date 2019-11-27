package cmd

import (
	"fmt"

	"golang.rgm.io/dwtk/internal/cli"
)

var InfoCmd = &cli.Command{
	Name:        "info",
	Description: "retrieve information from the target MCU and exit",
	Run: func(args []string) error {
		f, err := dw.ReadFuses()
		if err != nil {
			return err
		}

		fmt.Printf("Target MCU: %s\n", dw.MCU.Name)
		fmt.Printf("Fuses: low=0x%02X, high=0x%02X, extended=0x%02X, lockbit=0x%02X\n",
			f[0],
			f[1],
			f[2],
			f[3],
		)
		return nil
	},
}
