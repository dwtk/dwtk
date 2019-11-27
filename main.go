package main

import (
	"fmt"
	"os"

	"golang.rgm.io/dwtk/cmd"
)

func main() {
	cmd.Prog.AddCommand(cmd.DisableCmd)
	cmd.Prog.AddCommand(cmd.DumpCmd)
	cmd.Prog.AddCommand(cmd.EraseCmd)
	cmd.Prog.AddCommand(cmd.FlashCmd)
	cmd.Prog.AddCommand(cmd.GDBServerCmd)
	cmd.Prog.AddCommand(cmd.InfoCmd)
	cmd.Prog.AddCommand(cmd.ResetCmd)
	cmd.Prog.AddCommand(cmd.VerifyCmd)

	if err := cmd.Prog.Execute(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
