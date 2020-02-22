package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DisableCmd)
}

var DisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable debugWIRE in the target MCU and exit",
	Long:  "This command disables debugWIRE in the target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true
		return dw.Disable()
	},
}
