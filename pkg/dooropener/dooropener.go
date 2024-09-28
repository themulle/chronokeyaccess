package dooropener

import (
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type DoorOpener struct {
	RelayPin    int
	Chip        string
	ClosedState int
	OpenState   int
	Runtime     time.Duration
}

// NewWiegandReader erstellt einen neuen WiegandReader
func NewDoorOpener(relayPin int, closedState bool) *DoorOpener {
	opener := &DoorOpener{
		RelayPin: relayPin,
		Chip:     "gpiochip0",
		Runtime:  time.Second * 3,
	}
	if closedState {
		opener.ClosedState = 1
		opener.OpenState = 0
	} else {
		opener.ClosedState = 0
		opener.OpenState = 1
	}

	return opener
}

func (d *DoorOpener) InitAsOutput() error {
	line, err := gpiocdev.RequestLine(d.Chip, d.RelayPin, gpiocdev.AsOutput(d.ClosedState))
	if err != nil {
		return err
	}
	return line.Close()
}

func (d *DoorOpener) OpenDoor() error {
	line, err := gpiocdev.RequestLine(d.Chip, d.RelayPin, gpiocdev.AsOutput(d.ClosedState))
	if err != nil {
		return err
	}
	defer line.Close()
	defer line.SetValue(d.ClosedState)

	start := time.Now()
	for time.Since(start) < d.Runtime {
		if err = line.SetValue(d.OpenState); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 400)
		if err = line.SetValue(d.ClosedState); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}
