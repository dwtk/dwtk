package debugwire

func (dw *DebugWIRE) Disable() error {
	return dw.device.Write([]byte{0x06})
}

func (dw *DebugWIRE) Reset() error {
	if err := dw.SendBreak(); err != nil {
		return err
	}

	if err := dw.device.Write([]byte{0x07}); err != nil {
		return err
	}

	return dw.RecvBreak()
}

func (dw *DebugWIRE) GetSignature() (uint16, error) {
	if err := dw.device.Write([]byte{0xf3}); err != nil {
		return 0, err
	}

	return dw.device.ReadWord()
}
