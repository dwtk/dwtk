package cmd

import (
	"fmt"

	"golang.rgm.io/dwtk/internal/cli"
)

var DisableCmd = &cli.Command{
	Name:        "disable",
	Description: "disable debugWIRE in the target MCU, reset it and exit",
	Run: func(args []string) error {
		noReset = true

		if err := dw.Disable(); err != nil {
			return err
		}

		fmt.Println("Target device will be reseted and can be flashed using SPI.")
		fmt.Println("This must be done WITHOUT removing power from the device.")
		return nil
	},
}
