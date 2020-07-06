package cmd

import (
	"bytes"
	"fmt"

	"github.com/dwtk/dwtk/internal/hex"
	"github.com/spf13/cobra"
)

var (
	noVerifyEEPROM bool
)

func init() {
	EEPROMCmd.PersistentFlags().BoolVarP(
		&noVerifyEEPROM,
		"no-verify",
		"n",
		false,
		"do not verify data after writing",
	)

	RootCmd.AddCommand(EEPROMCmd)
}

var EEPROMCmd = &cobra.Command{
	Use:   "eeprom FILE",
	Short: "write data (Intel HEX) to target MCU's EEPROM, verify and exit",
	Long:  "This command writes data (Intel Hex) to target MCU's EEPROM, verifies and exits.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := hex.Parse(args[0])
		if err != nil {
			return err
		}

		cmd.Printf("Writing 0x%04x bytes to EEPROM ...\n", dw.MCU.EEPROMSize())
		if err := dw.WriteEEPROM(0, f); err != nil {
			return err
		}

		if noVerifyEEPROM {
			return nil
		}

		read := make([]byte, len(f))
		cmd.Printf("Verifying 0x%04x bytes from EEPROM ...\n", dw.MCU.EEPROMSize())
		if err := dw.ReadEEPROM(0, read); err != nil {
			return err
		}
		if bytes.Compare(f, read) != 0 {
			return fmt.Errorf("EEPROM mismatch: %v != %v", f, read)
		}

		return nil
	},
}
