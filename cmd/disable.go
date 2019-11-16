package cmd

import (
	"github.com/spf13/cobra"
)

var DisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable debugWIRE in the target MCU, reset it and exit",
	Long:  "This command disables debugWIRE in the target MCU, resets it and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true

		if err := dw.Disable(); err != nil {
			return err
		}

		cmd.Println("Target device will be reseted and can be flashed using SPI.")
		cmd.Println("This must be done WITHOUT removing power from the device.")
		return nil
	},
}
