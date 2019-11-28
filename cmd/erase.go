package cmd

import (
	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
)

var EraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "erase the target MCU's flash and exit",
	Long:  "This command erases the target MCU's flash and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true

		f, err := firmware.Empty(dw.MCU)
		if err != nil {
			return err
		}

		i := 1
		for _, page := range f.Pages {
			cmd.Printf("Erasing page %d/%d ...\n", i, len(f.Pages))
			if err := dw.WriteFlashPage(page.Address, page.Data); err != nil {
				return err
			}
			i += 1
		}

		return nil
	},
}
