package cmd

import (
	"github.com/dwtk/dwtk/avr"
	"github.com/spf13/cobra"
)

var (
	lfuse uint8
	hfuse uint8
	efuse uint8
	lock  uint8
)

func init() {
	FusesCmd.PersistentFlags().Uint8Var(
		&lfuse,
		"lfuse",
		0,
		"set low fuse",
	)
	FusesCmd.PersistentFlags().Uint8Var(
		&hfuse,
		"hfuse",
		0,
		"set high fuse",
	)
	FusesCmd.PersistentFlags().Uint8Var(
		&efuse,
		"efuse",
		0,
		"set extended fuse",
	)
	FusesCmd.PersistentFlags().Uint8Var(
		&lock,
		"lock",
		0,
		"set lock",
	)

	RootCmd.AddCommand(FusesCmd)
}

var FusesCmd = &cobra.Command{
	Use:   "fuses",
	Short: "retrieve or set fuses and lock from target MCU and exit",
	Long:  "This command retrieves or sets fuses and lock from target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		set := false

		if cmd.Flags().Changed("lfuse") {
			set = true
			if err := dw.WriteLFuse(lfuse); err != nil {
				return err
			}
		}
		if cmd.Flags().Changed("hfuse") {
			set = true
			if err := dw.WriteHFuse(hfuse); err != nil {
				return err
			}
		}
		if cmd.Flags().Changed("efuse") {
			set = true
			if err := dw.WriteEFuse(efuse); err != nil {
				return err
			}
		}
		if cmd.Flags().Changed("lock") {
			set = true
			if err := dw.WriteLock(lock); err != nil {
				return err
			}
		}

		if !set {
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
		}
		return nil
	},
}
