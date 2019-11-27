package cmd

import (
	"golang.rgm.io/dwtk/internal/cli"
)

var ResetCmd = &cli.Command{
	Name:        "reset",
	Description: "reset target MCU and exit",
	Run: func(args []string) error {
		// nothing to do, the post hook will reset.
		return nil
	},
}
