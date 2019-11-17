package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify FILE",
	Short: "verify firmware (ELF or Intel HEX) flashed to target MCU and exit",
	Long:  "This command verifies firmware (ELF or Intel Hex) flashed to target MCU and exits.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := firmware.Parse(args[0])
		if err != nil {
			return err
		}

		pages, err := dw.MCU.PrepareFirmware(f)
		if err != nil {
			return err
		}

		i := 1
		read := make([]byte, dw.MCU.FlashPageSize)
		for addr, data := range pages {
			cmd.Printf("Verifying page %d/%d ...\n", i, len(pages))
			if err := dw.ReadFlash(addr, read); err != nil {
				return err
			}
			if bytes.Compare(data, read) != 0 {
				return fmt.Errorf("Page mismatch 0x%04x: %v != %v", addr, data, read)
			}
			i += 1
		}

		return nil
	},
}
