package wieganddecoder

import (
	"fmt"
	"log/slog"
	"time"
)

type WiegandDecoder struct {
	MaxRuntime time.Duration
	MaxTimeout time.Duration
}

func NewWiegandDecoder() *WiegandDecoder {
	decoder := &WiegandDecoder{
		MaxRuntime: time.Millisecond * 10000,
		MaxTimeout: time.Millisecond * 100,
	}
	return decoder
}

func (d *WiegandDecoder) CheckInputComplete(bits []bool, start time.Time, stop time.Time) bool {
	if len(bits) == 0 {
		return false
	}

	if stop.Sub(start) > d.MaxRuntime || time.Since(start) > d.MaxRuntime {
		slog.Debug("runtime expired")
		return true
	}

	if len(bits) >= 34 || len(bits) >= 26 && time.Since(stop) > d.MaxTimeout {
		slog.Debug("no input for x seconds", "seconds", d.MaxTimeout.Seconds())
		return true
	}

	return false
}

func (d *WiegandDecoder) Decode(bits []bool) (uint64, error) {
	var dataBitsStart, dataBitsLength, parityCheckLength int
	var number uint64
	if len(bits) == 26 { // WG26 format
		dataBitsStart = 1      // Start after the first parity bit
		dataBitsLength = 24    // 24 bits of actual data
		parityCheckLength = 12 // First 12 and last 12 data bits for parity
	} else if len(bits) == 34 { // WG34 format
		dataBitsStart = 1      // Start after the first parity bit
		dataBitsLength = 32    // 32 bits of actual data
		parityCheckLength = 16 // First 16 and last 16 data bits for parity
	} else {
		return 0, fmt.Errorf("unsupported Wiegand format: bit slice length is %d", len(bits))
	}

	// Check even parity for WG26 (bits 1-12) or WG34 (bits 1-16)
	if !d.checkEvenParity(bits[1:dataBitsStart+parityCheckLength], bits[0]) {
		return 0, fmt.Errorf("even parity check failed")
	}

	// Check odd parity for WG26 (bits 13-24) or WG34 (bits 17-32)
	if !d.checkOddParity(bits[dataBitsStart+parityCheckLength:dataBitsStart+dataBitsLength], bits[len(bits)-1]) {
		return 0, fmt.Errorf("odd parity check failed")
	}

	for i := dataBitsStart; i <= dataBitsLength; i++ {
		if bits[i] {
			number |= 1 << (dataBitsLength - i) // Set the corresponding bit in the number
		}
	}
	return number, nil
}

// checkEvenParity checks if the number of 1's in the provided bits is even, compared to the parity bit.
func (d *WiegandDecoder) checkEvenParity(bits []bool, checkParity bool) bool {
	parity := false
	for _, bit := range bits {
		if bit {
			parity = !parity
		}
	}
	// Even parity: true if the number of 1's is even, false otherwise
	return parity == checkParity
}

// checkOddParity checks if the number of 1's in the provided bits is odd, compared to the parity bit.
func (d *WiegandDecoder) checkOddParity(bits []bool, checkParity bool) bool {
	parity := true
	for _, bit := range bits {
		if bit {
			parity = !parity
		}
	}
	// Odd parity: true if the number of 1's is odd, false otherwise
	return parity == checkParity
}
