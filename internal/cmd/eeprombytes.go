package cmd

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	noVerifyEEPROMBytes bool
	startEEPROMBytes    uint16
)

func init() {
	EEPROMBytesCmd.PersistentFlags().BoolVarP(
		&noVerifyEEPROMBytes,
		"no-verify",
		"n",
		false,
		"do not verify data after writing",
	)
	EEPROMBytesCmd.PersistentFlags().Uint16Var(
		&startEEPROMBytes,
		"start",
		0,
		"EEPROM address to start from",
	)

	RootCmd.AddCommand(EEPROMBytesCmd)
}

var EEPROMBytesCmd = &cobra.Command{
	Use:   "eeprom-bytes FILE",
	Short: "write arguments (as bytes) to target MCU's EEPROM, verify and exit",
	Long:  "This command writes arguments (as bytes) to target MCU's EEPROM, verifies and exits.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := []byte{}
		for _, arg := range args {
			b, err := strconv.ParseUint(arg, 0, 8)
			if err != nil {
				return err
			}
			f = append(f, byte(b))
		}

		cmd.Printf("Writing 0x%04x bytes to EEPROM, starting from 0x%04x ...\n", len(f), startEEPROMBytes)
		if err := dw.WriteEEPROM(startEEPROMBytes, f); err != nil {
			return err
		}

		if noVerifyEEPROMBytes {
			return nil
		}

		read := make([]byte, len(f))
		cmd.Printf("Verifying 0x%04x bytes from EEPROM, starting from 0x%04x ...\n", len(f), startEEPROMBytes)
		if err := dw.ReadEEPROM(startEEPROMBytes, read); err != nil {
			return err
		}
		if bytes.Compare(f, read) != 0 {
			return fmt.Errorf("EEPROM mismatch: %v != %v", f, read)
		}

		return nil
	},
}
