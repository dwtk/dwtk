package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EraseChipCmd)
}

var EraseChipCmd = &cobra.Command{
	Use:   "erase-chip",
	Short: "erase target flash, eeprom, lock using SPI command and exit",
	Long:  "This command erases target flash, eeprom, lock using SPI command and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true
		return dw.ChipErase()
	},
}
