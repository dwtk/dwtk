package main

import (
	"os"

	"golang.rgm.io/dwtk/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
