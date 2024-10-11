package dooropener

import (
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type GpioOut struct {
	Pin    		int
	Chip        string
	OnState		int
	OffState	int
}

func NewGpioOut(relayPin int, offState bool) *GpioOut {
	retval := &GpioOut{
		Pin: relayPin,
		Chip:     "gpiochip0",
	}
	if offState {
		retval.OnState = 1
		retval.OffState = 0
	} else {
		retval.OnState = 0
		retval.OffState = 1
	}

	return retval
}

func (g *GpioOut) InitAsOutput() error {
	line, err := gpiocdev.RequestLine(g.Chip, g.Pin, gpiocdev.AsOutput(g.OffState))
	if err != nil {
		return err
	}
	return line.Close()
}

func (g *GpioOut) ActivateFor(runtime time.Duration) error {
	line, err := gpiocdev.RequestLine(g.Chip, g.Pin, gpiocdev.AsOutput(g.OffState))
	if err != nil {
		return err
	}
	defer line.Close()
	defer line.SetValue(g.OffState)

	if err = line.SetValue(g.OnState); err != nil {
		return err
	}
	time.Sleep(runtime)
	if err = line.SetValue(g.OffState); err != nil {
		return err
	}
	return nil
}
