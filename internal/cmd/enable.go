package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EnableCmd)
}

var EnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable debugWIRE in the target MCU and exit",
	Long:  "This command enables debugWIRE in the target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true
		return dw.Enable()
	},
}
