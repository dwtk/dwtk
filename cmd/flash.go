package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
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

		i := 1
		for _, page := range pages {
			cmd.Printf("Flashing page %d/%d ...\n", i, len(pages))
			if err := dw.WriteFlashPage(page.Address, page.Data); err != nil {
				return err
			}
			i += 1
		}

		if noVerify {
			return nil
		}

		i = 1
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
