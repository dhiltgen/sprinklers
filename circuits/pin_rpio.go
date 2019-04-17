package circuits

import (
	"github.com/stianeikeland/go-rpio"
)

func newGPIOPin(pin uint8) Pin {
	return rpio.Pin(pin)
}

func gpioPinInit() error {
	return rpio.Open()
}

func init() {
	newPin = newGPIOPin
	pinsInit = gpioPinInit
}
