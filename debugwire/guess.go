package debugwire

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.rgm.io/dwtk/logger"
	"golang.rgm.io/dwtk/usbserial"
)

func GuessDevice() (string, error) {
	matches, err := filepath.Glob("/dev/ttyUSB*")
	if err != nil {
		return "", err
	}
	if matches == nil {
		return "", fmt.Errorf("usbserial: no USB serial port found")
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("usbserial: more than one USB serial port found: %s",
			strings.Join(matches, ", "))
	}

	logger.Debug.Printf(" * Detected serial port: %s", matches[0])
	return matches[0], nil
}

func GuessBaudrate(portDevice string) (uint32, error) {
	// max supported mcu frequency is 20MHz, there are faster AVRs, but they
	// provide better capabilities than debugwire. we want the faster baudrate
	// possible
	for i := 20; i > 0; i-- {
		baudrate := uint32((i * 1000000) / 128)
		p, err := usbserial.Open(portDevice, baudrate)
		if err != nil {
			return 0, err
		}

		if err := p.SendBreak(); err != nil {
			return 0, err
		}

		c, err := p.RecvBreak()
		if err != nil {
			return 0, err
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

	return 0, fmt.Errorf("debugwire: failed to detect baudrate for serial port: %s", portDevice)
}
