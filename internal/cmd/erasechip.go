package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EraseChipCmd)
}

var EraseChipCmd = &cobra.Command{
	Use:   "erase-chip",
	Short: "erase target chip (flash, eeprom, lockbits) and exit",
	Long:  "This command erases the target chip (flash, eeprom, lockbits) and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true
		return dw.ChipErase()
	},
}
