package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
)

func init() {
	RootCmd.AddCommand(VerifyCmd)
}

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

		read := make([]byte, dw.MCU.FlashPageSize)
		for i, page := range pages {
			cmd.Printf("Verifying page 0x%04x (%d/%d) ...\n", page.Address, i+1, len(pages))
			if err := dw.ReadFlash(page.Address, read); err != nil {
				return err
			}
			if bytes.Compare(page.Data, read) != 0 {
				return fmt.Errorf("Page mismatch 0x%04x: %v != %v", page.Address, page.Data, read)
			}
		}

		return nil
	},
}
