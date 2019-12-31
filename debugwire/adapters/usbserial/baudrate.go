package usbserial

import (
	"fmt"

	"github.com/dwtk/dwtk/internal/usbserial"
)

func detectBaudrate(serialPort string) (uint32, error) {
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
			return baudrate, nil
		}
	}

	return 0, fmt.Errorf("debugwire: usbserial: failed to detect baudrate for serial port: %s", serialPort)
}
