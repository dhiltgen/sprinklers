package circuits

// NewPinFunc creates a new pin
type NewPinFunc func(gpio uint8) Pin

// PinInitFunc initializes a pin
type PinInitFunc func() error

var newPin NewPinFunc
var pinsInit PinInitFunc

// Pin is an abstraction on top of go-rpio to facilitate unit tests
type Pin interface {
	Output()
	High()
	Low()
}
