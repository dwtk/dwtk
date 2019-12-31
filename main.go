package main

import (
	"os"

	"github.com/dwtk/dwtk/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
