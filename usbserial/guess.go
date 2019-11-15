package usbserial

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.rgm.io/dwtk/logger"
)

func GuessPortDevice() (string, error) {
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
