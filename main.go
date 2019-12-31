package main

import (
	"os"

	"github.com/dwtk/dwtk/internal/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
