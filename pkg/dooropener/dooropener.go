package dooropener

import (
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type DoorOpener struct {
	GpioOut
}

func NewDoorOpener(relayPin int, offState bool) *DoorOpener {
	retval := &DoorOpener{}
	retval.GpioOut=*NewGpioOut(relayPin,offState)
	return retval
}

func (g *DoorOpener) ActivateFor(runtime time.Duration)  error {
	line, err := gpiocdev.RequestLine(g.Chip, g.Pin, gpiocdev.AsOutput(g.OffState))
	if err != nil {
		return err
	}
	defer line.Close()
	defer line.SetValue(g.OffState)

	start := time.Now()
	for time.Since(start) < runtime {
		if err = line.SetValue(g.OnState); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 400)
		if err = line.SetValue(g.OffState); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}
