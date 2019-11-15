package cmd

import (
	"github.com/spf13/cobra"
	"golang.rgm.io/dwtk/debugwire"
	"golang.rgm.io/dwtk/logger"
	"golang.rgm.io/dwtk/usbserial"
)

var (
	dw      *debugwire.DebugWire
	noReset bool

	serialPort string
	baudrate   uint32
	debug      bool

	Version = "git"
)

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&serialPort,
		"serial-port",
		"s",
		"",
		"Serial port device (e.g. /dev/ttyUSB0. Default: detect)",
	)
	RootCmd.PersistentFlags().Uint32VarP(
		&baudrate,
		"baudrate",
		"b",
		0,
		"Serial communication baudrate (e.g. 62500. Default: detect)",
	)
	RootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"Enable debugging messages (Default: false)",
	)
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

var RootCmd = &cobra.Command{
	Use:          "dw",
	Short:        "debugWire toolkit",
	Long:         "debugWire toolkit",
	Version:      Version,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			logger.EnableDebug()
		}

		if serialPort == "" {
			var err error
			serialPort, err = usbserial.GuessPortDevice()
			if err != nil {
				return err
			}
		}

		if baudrate == 0 {
			var err error
			baudrate, err = debugwire.GuessBaudrate(serialPort)
			if err != nil {
				return err
			}
		}

		var err error
		dw, err = debugwire.New(serialPort, baudrate)
		if err != nil {
			return err
		}

		noReset = false

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if dw != nil {
			defer dw.Close()
			if !noReset {
				if err := dw.Reset(); err != nil {
					return err
				}
				if err := dw.Go(); err != nil {
					return err
				}
			}
		}

		return nil
	},
}
