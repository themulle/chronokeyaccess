package wiegand

import (
	"fmt"
	"log/slog"
	"sync"
	"syscall"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type WiegandReader struct {
	Data0Pin int
	Data1Pin int
	Chip     string

	FirstInput time.Time
	LastInput  time.Time

	bits         []bool
	mutex        sync.Mutex
	requestLines *gpiocdev.Lines
}

// NewWiegandReader erstellt einen neuen WiegandReader
func NewWiegandReader(data0Pin, data1Pin int) *WiegandReader {
	reader := &WiegandReader{
		Data0Pin: data0Pin,
		Data1Pin: data1Pin,
		Chip:     "gpiochip0",
		bits:     make([]bool, 0),
	}
	return reader
}

func (w *WiegandReader) Start() error {
	var err error
	w.requestLines, err = gpiocdev.RequestLines(w.Chip, []int{w.Data0Pin, w.Data1Pin},
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithEventHandler(w.processGpioEvent))
	if err != nil {
		if err == syscall.Errno(22) {
			err = fmt.Errorf("Note that the WithPullUp option requires Linux 5.5 or later - check your kernel version: %w", err)
		}
		return fmt.Errorf("RequestLine returned error: %w", err)
	}
	return err
}

func (w *WiegandReader) Close() {
	if w.requestLines != nil {
		w.requestLines.Close()
		w.requestLines = nil
	}
}

func (w *WiegandReader) GetCache() ([]bool, time.Time, time.Time) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	bits := w.bits
	return bits, w.FirstInput, w.LastInput
}

func (w *WiegandReader) ResetCache() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.bits = make([]bool, 0)
	w.FirstInput = time.Time{}
	w.LastInput = time.Time{}
}

func (w *WiegandReader) processGpioEvent(evt gpiocdev.LineEvent) {
	slog.Debug("gpio event", "offset", evt.Offset, "eventtype", evt.Type)
	if evt.Type == gpiocdev.LineEventFallingEdge {
		if w.FirstInput.IsZero() {
			w.FirstInput = time.Now()
		}
		w.LastInput = time.Now()
		if evt.Offset == w.Data0Pin {
			w.bits = append(w.bits, false)
		} else if evt.Offset == w.Data1Pin {
			w.bits = append(w.bits, true)
		}
	}
}
