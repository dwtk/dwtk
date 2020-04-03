package cmd

import (
	"github.com/dwtk/dwtk/internal/hex"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DumpEEPROMCmd)
}

var DumpEEPROMCmd = &cobra.Command{
	Use:   "dump-eeprom FILE",
	Short: "dump data (Intel HEX) from target MCU' EEPROM and exit",
	Long:  "This command dumps data (Intel Hex) from target MCU's EEPROM and exits.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		read := make([]byte, dw.MCU.EEPROMSize)
		cmd.Printf("Retrieving 0x%04x bytes from EEPROM ...\n", dw.MCU.EEPROMSize)
		if err := dw.ReadEEPROM(0, read); err != nil {
			return err
		}
		return hex.Dump(args[0], read)
	},
}
