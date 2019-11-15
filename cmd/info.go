package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve information from the target MCU and exit",
	Long:  "This command retrieves information from the target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dw.MCU == nil {
			return fmt.Errorf("Failed to detect MCU")
		}

		f, err := dw.ReadFuses()
		if err != nil {
			return err
		}

		cmd.Printf("Target MCU: %s\n", dw.MCU.Name)
		cmd.Printf("Fuses: low=0x%02X, high=0x%02X, extended=0x%02X, lockbit=0x%02X\n",
			f[0],
			f[1],
			f[2],
			f[3],
		)
		return nil
	},
}
