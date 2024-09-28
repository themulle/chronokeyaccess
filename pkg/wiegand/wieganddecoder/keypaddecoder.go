package wieganddecoder

import (
	"fmt"
	"log/slog"
	"time"
)

type KeypadDecoder struct {
	MaxRuntime time.Duration
	MaxTimeout time.Duration
}

func NewKeypadDecoder() *KeypadDecoder {
	decoder := &KeypadDecoder{
		MaxRuntime: time.Second * 10,
		MaxTimeout: time.Second * 3,
	}
	return decoder
}

func (d *KeypadDecoder) CheckInputComplete(bits []bool, start time.Time, stop time.Time) bool {
	if len(bits) == 0 {
		return false
	}
	if stop.Sub(start) > d.MaxRuntime || time.Since(start) > d.MaxRuntime {
		slog.Debug("runtime expired")
		return true
	}

	if time.Since(stop) > d.MaxTimeout {
		slog.Debug("no input for x seconds", "seconds", d.MaxTimeout.Seconds())
		return true
	}

	//very short inputs are provided by card
	if time.Since(start) < time.Second {
		slog.Debug("time too short")
		return false
	}

	if len(bits) >= 8 && len(bits)%4 == 0 {
		var lastNumber uint64
		last4Bits := bits[len(bits)-4:]
		for i := 0; i < 4; i++ {
			if last4Bits[i] {
				lastNumber |= 1 << (3 - i) // Shift the bit to the correct position
			}
		}
		if lastNumber > 9 {
			slog.Debug("KeypadDecoder magic number detected", "number", lastNumber)
			return true
		}
	}

	return false
}

func (d *KeypadDecoder) Decode(bits []bool) (uint64, error) {
	bitLength := len(bits)
	if bitLength == 0 || bitLength%4 != 0 {
		return 0, fmt.Errorf("unsupported Keypad format: bit slice length is %d", len(bits))
	}

	var number uint64
	for i := 0; i < bitLength; i += 4 {
		// Extract the 4 bits
		var digit uint64
		for j := 0; j < 4; j++ {
			if bits[i+j] {
				digit |= 1 << (3 - j) // Build the 4-bit digit
			}
		}
		// If the digit is greater than 9, it is invalid for BCD encoding
		if digit > 9 {
			//last bit may be greater than 9 to signale "enter"
			if i+4 != bitLength {
				return 0, fmt.Errorf("invalid BCD encoding at position %d", i/4)
			}
		} else {
			number = number*10 + digit
		}

	}

	return number, nil
}
