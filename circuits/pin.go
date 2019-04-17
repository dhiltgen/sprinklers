package circuits

// Abstraction on top of go-rpio to facilitate unit tests

type NewPinFunc func(gpio uint8) Pin
type PinInitFunc func() error

var newPin NewPinFunc
var pinsInit PinInitFunc

type Pin interface {
	Output()
	High()
	Low()
}
