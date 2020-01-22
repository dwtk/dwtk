package dwtkice

func (dw *DwtkIceAdapter) ReadFuses() ([]byte, error) {
	f := make([]byte, 4)
	if err := dw.controlIn(cmdReadFuses, 0, 0, f); err != nil {
		return nil, err
	}
	return f, nil
}
