package cmd

import (
	"github.com/dwtk/dwtk/gdbserver"
	"github.com/spf13/cobra"
)

var (
	addr          string
	disableTimers bool
)

func init() {
	GDBServerCmd.PersistentFlags().StringVarP(
		&addr,
		"addr",
		"a",
		"localhost:8000",
		"GDB server host:port",
	)
	GDBServerCmd.PersistentFlags().BoolVarP(
		&disableTimers,
		"disable-timers",
		"t",
		false,
		"disable timers",
	)

	RootCmd.AddCommand(GDBServerCmd)
}

var GDBServerCmd = &cobra.Command{
	Use:   "gdbserver",
	Short: "start remote debugging session for GDB",
	Long:  "This command starts a remote debuggins session for GDB.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dw.Timers = !disableTimers
		return gdbserver.ListenAndServe(addr, dw)
	},
}
