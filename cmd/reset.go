package cmd

import (
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset target MCU and exit",
	Long:  "This command resets the target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// nothing to do, the post hook will reset.
		return nil
	},
}
