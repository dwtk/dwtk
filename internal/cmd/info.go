package cmd

import (
	"github.com/dwtk/dwtk/avr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(InfoCmd)
}

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "retrieve device information from target MCU and exit",
	Long:  "This command retrieves device information from target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Print(dw.Info())
		cmd.Printf("\n")

		cmd.Printf("Target MCU: %s\n", dw.MCU.Name())

		f, err := dw.ReadFuses()
		if err != nil {
			return err
		}
		cmd.Printf("Fuses: low=0x%02X, high=0x%02X, extended=0x%02X, lockbit=0x%02X\n",
			f[avr.LOW_FUSE],
			f[avr.HIGH_FUSE],
			f[avr.EXTENDED_FUSE],
			f[avr.LOCKBIT],
		)
		return nil
	},
}
