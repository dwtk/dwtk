package cmd

import (
	"bytes"
	"fmt"

	"github.com/dwtk/dwtk/firmware"
	"github.com/spf13/cobra"
)

var (
	noVerify bool
)

func init() {
	FlashCmd.PersistentFlags().BoolVarP(
		&noVerify,
		"no-verify",
		"n",
		false,
		"do not verify flashed firmware",
	)

	RootCmd.AddCommand(FlashCmd)
}

var FlashCmd = &cobra.Command{
	Use:   "flash FILE",
	Short: "flash firmware (ELF or Intel HEX) to target MCU, verify and exit",
	Long:  "This command flashes firmware (ELF or Intel Hex) to target MCU, verifies and exits.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := firmware.NewFromFile(args[0], dw.MCU)
		if err != nil {
			return err
		}

		pages := f.SplitPages()

		for i, page := range pages {
			cmd.Printf("Flashing page 0x%04x (%d/%d) ...\n", page.Address, i+1, len(pages))
			if err := dw.WriteFlashPage(page.Address, page.Data); err != nil {
				return err
			}
		}

		if noVerify {
			return nil
		}

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
