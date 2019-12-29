package dwtk

import (
	"fmt"
)

func (dw *DwtkAdapter) ubrrToBaudrate(ubrr uint16) (uint32, error) {
	presc, freq, err := dw.getBaudratePrescaler()
	if err != nil {
		return 0, err
	}

	if freq == 0 {
		return 0, fmt.Errorf("debugwire: dwtk: got invalid oscillator frequency")
	}

	return (uint32(freq) * 1000000) / uint32(uint16(presc)*(ubrr+1)), nil
}

func (dw *DwtkAdapter) baudrateToUbrr(baudrate uint32) (uint16, error) {
	presc, freq, err := dw.getBaudratePrescaler()
	if err != nil {
		return 0, err
	}

	if freq == 0 {
		return 0, fmt.Errorf("debugwire: dwtk: got invalid oscillator frequency")
	}

	return uint16((uint32(freq)*1000000)/(uint32(presc)*baudrate)) - 1, nil
}

func (dw *DwtkAdapter) getBaudratePrescaler() (byte, byte, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetBaudratePrescaler, 0, 0, f); err != nil {
		return 0, 0, err
	}
	return f[0], f[1], nil
}

func (dw *DwtkAdapter) detectBaudrate() (uint16, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdDetectBaudrate, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkAdapter) setBaudrate(baudrate uint16) error {
	return dw.controlIn(cmdSetBaudrate, baudrate, 0, nil)
}
