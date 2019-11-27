package cmd

import (
	"flag"

	"golang.rgm.io/dwtk/gdbserver"
	"golang.rgm.io/dwtk/internal/cli"
)

var (
	addr          string
	disableTimers bool
)

var GDBServerCmd = &cli.Command{
	Name:        "gdbserver",
	Description: "start remote debugging session for GDB",
	SetFlags: func(fs *flag.FlagSet) error {
		fs.StringVar(
			&addr,
			"a",
			"localhost:8000",
			"GDB server host:port",
		)
		fs.BoolVar(
			&disableTimers,
			"t",
			false,
			"disable timers",
		)
		return nil
	},
	Run: func(args []string) error {
		dw.Timers = !disableTimers
		return gdbserver.ListenAndServe(addr, dw)
	},
}
