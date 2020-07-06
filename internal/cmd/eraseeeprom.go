package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EraseEEPROMCmd)
}

var EraseEEPROMCmd = &cobra.Command{
	Use:   "erase-eeprom",
	Short: "erase target MCU's EEPROM and exit",
	Long:  "This command erases target MCU's EEPROM and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := make([]byte, dw.MCU.EEPROMSize())
		for i := uint16(0); i < dw.MCU.EEPROMSize(); i++ {
			f[i] = 0xff
		}

		cmd.Printf("Erasing 0x%04x bytes from EEPROM ...\n", dw.MCU.EEPROMSize())
		return dw.WriteEEPROM(0, f)
	},
}
