package circuits

type dummyPin struct {
	Level bool
}

func (p dummyPin) Output() {}
func (p dummyPin) High() {
	p.Level = true
}
func (p dummyPin) Low() {
	p.Level = false
}

func newDummyPin(_ uint8) Pin {
	return dummyPin{}
}

func DummyInit() {
	newPin = newDummyPin
	pinsInit = func() error { return nil }
}
