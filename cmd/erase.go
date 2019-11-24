package cmd

import (
	"github.com/spf13/cobra"
)

var EraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "erase the target MCU's flash and exit",
	Long:  "This command erases the target MCU's flash and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		noReset = true

		pages, err := dw.MCU.PrepareFirmware(make([]byte, dw.MCU.FlashSize))
		if err != nil {
			return err
		}

		i := 1
		for addr, data := range pages {
			cmd.Printf("Erasing page %d/%d ...\n", i, len(pages))
			if err := dw.WriteFlashPage(addr, data); err != nil {
				return err
			}
			i += 1
		}

		return nil
	},
}
