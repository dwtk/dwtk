package cmd

import (
	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/firmware"
)

func init() {
	RootCmd.AddCommand(EraseCmd)
}

var EraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "erase the target MCU's flash and exit",
	Long:  "This command erases the target MCU's flash and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true

		f, err := firmware.NewEmpty(dw.MCU)
		if err != nil {
			return err
		}

		pages := f.SplitPages()

		i := 1
		for _, page := range pages {
			cmd.Printf("Erasing page 0x%04x (%d/%d) ...\n", page.Address, i, len(pages))
			if err := dw.WriteFlashPage(page.Address, page.Data); err != nil {
				return err
			}
			i += 1
		}

		return nil
	},
}
