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
		numPages := dw.MCU.FlashSize / dw.MCU.FlashPageSize

		read := make([]byte, dw.MCU.FlashPageSize)
		f := []byte{}
		for i := uint16(0); i < numPages; i += 1 {
			addr := i * dw.MCU.FlashPageSize
			cmd.Printf("Retrieving page 0x%04x (%d/%d) ...\n", addr, i+1, numPages)
			if err := dw.ReadFlash(addr, read); err != nil {
				return err
			}
			f = append(f, read...)
		}

		fw, err := firmware.NewFromData(f, dw.MCU)
		if err != nil {
			return err
		}

		return fw.Dump(args[0])
	},
}
