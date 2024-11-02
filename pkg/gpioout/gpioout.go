package gpioout

import (
	"log/slog"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type GpioOut struct {
	Pin    		int
	Chip        string
	DefaultState	int
}

func NewGpioOut(relayPin int, defaultState int) *GpioOut {
	retval := &GpioOut{
		Pin: relayPin,
		Chip:     "gpiochip0",
		DefaultState: defaultState,
	}

	return retval
}

func (g *GpioOut) InitAsOutput() error {
	line, err := gpiocdev.RequestLine(g.Chip, g.Pin, gpiocdev.AsOutput(g.DefaultState))
	if err != nil {
		return err
	}
	return line.Close()
}

func (g *GpioOut) ActivateFor(runtime time.Duration) error {
	line, err := gpiocdev.RequestLine(g.Chip, g.Pin, gpiocdev.AsOutput(g.DefaultState))
	if err != nil {
		return err
	}
	defer line.Close()
	defer line.SetValue(g.DefaultState)

	onValue:=1
	if g.DefaultState>0 {
		onValue=0
	}
	slog.Debug("setting pin", "pin", g.Pin, "value", onValue)
	if err = line.SetValue(onValue); err != nil {
		return err
	}
	time.Sleep(runtime)
	slog.Debug("setting pin", "pin", g.Pin, "value", g.DefaultState)
	if err = line.SetValue(g.DefaultState); err != nil {
		return err
	}
	return nil
}
