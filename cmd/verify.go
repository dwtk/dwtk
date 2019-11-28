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
		f, err := firmware.NewFromFile(args[0], dw.MCU)
		if err != nil {
			return err
		}

		pages := f.SplitPages()

		i := 1
		read := make([]byte, dw.MCU.FlashPageSize)
		for _, page := range pages {
			cmd.Printf("Verifying page %d/%d ...\n", i, len(pages))
			if err := dw.ReadFlash(page.Address, read); err != nil {
				return err
			}
			if bytes.Compare(page.Data, read) != 0 {
				return fmt.Errorf("Page mismatch 0x%04x: %v != %v", page.Address, page.Data, read)
			}
			i += 1
		}

		return nil
	},
}
