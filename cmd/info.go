package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(InfoCmd)
}

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "retrieve information from the target MCU and exit",
	Long:  "This command retrieves information from the target MCU and exits.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := dw.ReadFuses()
		if err != nil {
			return err
		}

		cmd.Printf("Serial Port: %s\n", serialPort)
		cmd.Printf("Baud Rate: %d\n", baudrate)
		cmd.Printf("\n")
		cmd.Printf("Target MCU: %s\n", dw.MCU.Name)
		cmd.Printf("Fuses: low=0x%02X, high=0x%02X, extended=0x%02X, lockbit=0x%02X\n",
			f[0],
			f[1],
			f[2],
			f[3],
		)
		return nil
	},
}
