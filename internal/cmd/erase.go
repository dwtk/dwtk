package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EraseCmd)
}

var EraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "erase target MCU's flash and exit",
	Long:  "This command erases target MCU's flash and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true

		numPages := dw.MCU.FlashSize / uint16(dw.MCU.FlashPageSize)
		for i := uint16(0); i < numPages; i++ {
			address := i * dw.MCU.FlashPageSize
			cmd.Printf("Erasing page 0x%04x (%d/%d) ...\n", address, i+1, numPages)
			if err := dw.EraseFlashPage(address); err != nil {
				return err
			}
		}

		return nil
	},
}
