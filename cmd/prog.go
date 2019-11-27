package cmd

import (
	"flag"
	"fmt"

	"golang.rgm.io/dwtk/debugwire"
	"golang.rgm.io/dwtk/internal/cli"
	"golang.rgm.io/dwtk/logger"
	"golang.rgm.io/dwtk/version"
)

var (
	dw      *debugwire.DebugWire
	noReset bool

	serialPort string
	baudrate   uint
	debug      bool
)

var Prog = &cli.Program{
	Version:     version.Version,
	Description: "debugWIRE toolkit for AVR microcontrollers",
	SetFlags: func(fs *flag.FlagSet) error {
		fs.StringVar(
			&serialPort,
			"s",
			"",
			"serial port device (e.g. /dev/ttyUSB0) (Default: detect)",
		)
		fs.UintVar(
			&baudrate,
			"b",
			0,
			"serial port baudrate (e.g. 62500) (Default: detect)",
		)
		fs.BoolVar(
			&debug,
			"d",
			false,
			"enable debugging messages",
		)
		return nil
	},
	Pre: func(args []string) error {
		if debug {
			logger.EnableDebug()
		}

		if serialPort == "" {
			var err error
			serialPort, err = debugwire.GuessDevice()
			if err != nil {
				return err
			}
		}

		baudrate := uint32(baudrate)
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

		if dw.MCU == nil {
			return fmt.Errorf("Failed to detect MCU")
		}

		noReset = false

		return nil
	},
	Post: func(args []string) error {
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
