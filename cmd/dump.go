package cmd

import (
	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
)

var DumpCmd = &cobra.Command{
	Use:   "dump FILE",
	Short: "dump firmware (Intel HEX) from target MCU and exit",
	Long:  "This command dumps firmware (Intel Hex) from target MCU and exits.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		numPages, err := dw.MCU.NumFlashPages()
		if err != nil {
			return err
		}

		read := make([]byte, dw.MCU.FlashPageSize)
		f := []byte{}
		for i := uint16(0); i < numPages; i += 1 {
			cmd.Printf("Retrieving page %d/%d ...\n", i+1, numPages)
			addr := i * dw.MCU.FlashPageSize
			if err := dw.ReadFlash(addr, read); err != nil {
				return err
			}
			f = append(f, read...)
		}

		return firmware.Dump(args[0], f)
	},
}
