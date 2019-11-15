package debugwire

func (dw *DebugWire) Disable() error {
	return dw.Port.Write([]byte{0x06})
}

func (dw *DebugWire) Reset() error {
	if err := dw.SendBreak(); err != nil {
		return err
	}

	if err := dw.Port.Write([]byte{0x07}); err != nil {
		return err
	}

	return dw.RecvBreak()
}

func (dw *DebugWire) GetSignature() (uint16, error) {
	if err := dw.Port.Write([]byte{0xf3}); err != nil {
		return 0, err
	}

	return dw.Port.ReadWord()
}
