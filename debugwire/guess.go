package debugwire

import (
	"fmt"

	"golang.rgm.io/dwtk/logger"
	"golang.rgm.io/dwtk/usbserial"
)

func GuessBaudrate(portDevice string) (uint32, error) {
	// max supported mcu frequency is 20MHz, there are faster AVRs, but they
	// provide better capabilities than debugwire. we want the faster baudrate
	// possible
	for i := 20; i > 0; i-- {
		baudrate := uint32((i * 1000000) / 128)
		fd, err := usbserial.Open(portDevice, baudrate)
		if err != nil {
			return 0, err
		}

		c, err := usbserial.SendBreak(fd)
		if err != nil {
			return 0, err
		}

		// we could reuse the fd, but we won't, for reproducibility reasons
		if err := usbserial.Close(fd); err != nil {
			return 0, err
		}

		if c == 0x55 {
			logger.Debug.Printf(" * Detected baudrate: %d\n", baudrate)
			return baudrate, nil
		}
	}

	return 0, fmt.Errorf("debugwire: failed to detect baudrate for serial port: %s", portDevice)
}
