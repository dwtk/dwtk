package dwtkice

import (
	"fmt"
)

func (dw *DwtkIceAdapter) ubrrToBaudrate(ubrr uint16) (uint32, error) {
	presc, freq, err := dw.getBaudratePrescaler()
	if err != nil {
		return 0, err
	}

	if freq == 0 {
		return 0, fmt.Errorf("debugwire: dwtk-ice: got invalid oscillator frequency")
	}

	return (uint32(freq) * 1000000) / uint32(uint16(presc)*(ubrr+1)), nil
}

func (dw *DwtkIceAdapter) baudrateToUbrr(baudrate uint32) (uint16, error) {
	presc, freq, err := dw.getBaudratePrescaler()
	if err != nil {
		return 0, err
	}

	if freq == 0 {
		return 0, fmt.Errorf("debugwire: dwtk-ice: got invalid oscillator frequency")
	}

	tmp := uint32(((uint64(freq) * 10000000) / (uint64(presc) * uint64(baudrate))) - 10)
	ubrr := uint16(tmp / 10)
	if tmp%10 >= 5 {
		ubrr++
	}
	return ubrr, nil
}

func (dw *DwtkIceAdapter) getBaudratePrescaler() (byte, byte, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdGetBaudratePrescaler, 0, 0, f); err != nil {
		return 0, 0, err
	}
	return f[0], f[1], nil
}

func (dw *DwtkIceAdapter) detectBaudrate() (uint16, error) {
	f := make([]byte, 2)
	if err := dw.controlIn(cmdDetectBaudrate, 0, 0, f); err != nil {
		return 0, err
	}
	return (uint16(f[0]) << 8) | uint16(f[1]), nil
}

func (dw *DwtkIceAdapter) setBaudrate(baudrate uint16) error {
	return dw.controlIn(cmdSetBaudrate, baudrate, 0, nil)
}
