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
		"Do not verify flashed firmware",
	)
}

var FlashCmd = &cobra.Command{
	Use:   "flash FILE",
	Short: "Flash Intel Hex program to target MCU, verify and exit",
	Long:  "This command flashes Intel Hex program to target MCU, verifies and exits.",
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
		for addr, data := range pages {
			cmd.Printf("Flashing page %d/%d ...\n", i, len(pages))
			if err := dw.WriteFlashPage(addr, data); err != nil {
				return err
			}
			i += 1
		}

		if noVerify {
			return nil
		}

		i = 1
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
