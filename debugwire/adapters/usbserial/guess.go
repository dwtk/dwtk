package usbserial

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.rgm.io/dwtk/logger"
	"golang.rgm.io/dwtk/usbserial"
)

func guessSerialPort() (string, error) {
	matches, err := filepath.Glob("/dev/ttyUSB*")
	if err != nil {
		return "", err
	}
	if matches == nil {
		return "", fmt.Errorf("debugwire: usbserial: no USB serial port found")
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("debugwire: usbserial: more than one USB serial port found: %s",
			strings.Join(matches, ", "))
	}

	logger.Debug.Printf(" * Detected serial port: %s", matches[0])
	return matches[0], nil
}

func guessBaudrate(serialPort string) (uint32, error) {
	// max supported mcu frequency is 20MHz, there are faster AVRs, but they
	// provide better capabilities than debugwire. we want the faster baudrate
	// possible
	for i := 20; i > 0; i-- {
		baudrate := uint32((i * 1000000) / 128)
		p, err := usbserial.Open(serialPort, baudrate)
		if err != nil {
			return 0, err
		}

		if err := p.SendBreak(); err != nil {
			return 0, err
		}

		// if devices are running on very low frequency (e.g. internal clock with
		// CKDIV8 enabled) we may not even get a break response.
		c, err := p.RecvBreak()
		if err != nil {
			if err := p.Close(); err != nil {
				return 0, err
			}
			continue
		}

		// we could reuse the instance, but we won't, for reproducibility reasons
		if err := p.Close(); err != nil {
			return 0, err
		}

		if c == 0x55 {
			logger.Debug.Printf(" * Detected baudrate: %d\n", baudrate)
			return baudrate, nil
		}
	}

	return 0, fmt.Errorf("debugwire: usbserial: failed to detect baudrate for serial port: %s", serialPort)
}
