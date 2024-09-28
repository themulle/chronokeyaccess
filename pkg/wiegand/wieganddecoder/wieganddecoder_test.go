package wieganddecoder_test

import (
	"testing"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/wiegand/wieganddecoder"
)

// Helper function to convert a binary string to a bool slice for easier test case setup

func TestWiegandDecoder_CheckInputComplete(t *testing.T) {
	decoder := wieganddecoder.NewWiegandDecoder()

	t.Run("No bits provided", func(t *testing.T) {
		start := time.Now()
		stop := start.Add(time.Millisecond * 100)
		bits := []bool{}
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("Runtime exceeded", func(t *testing.T) {
		start := time.Now().Add(-time.Second * 15)
		stop := time.Now()
		bits := bitsFromString("110010")
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (runtime exceeded), got %v", result)
		}
	})

	t.Run("Timeout exceeded with WG26", func(t *testing.T) {
		start := time.Now().Add(-time.Second * 5)
		stop := time.Now().Add(-time.Second * 5)
		bits := bitsFromString("10000000000101111010110011") // 26 bits
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (timeout exceeded), got %v", result)
		}
	})

	t.Run("Valid input with WG34", func(t *testing.T) {
		start := time.Now().Add(-time.Second * 2)
		stop := time.Now().Add(-time.Millisecond * 500)
		bits := bitsFromString("1000000000010111101011001000110001") // 34 bits
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (valid WG34 input), got %v", result)
		}
	})

	t.Run("Input too short for WG26", func(t *testing.T) {
		start := time.Now().Add(-time.Millisecond * 500)
		stop := time.Now()
		bits := bitsFromString("100110") // too short
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != false {
			t.Errorf("Expected false (input too short), got %v", result)
		}
	})
}

func TestWiegandDecoder_Decode(t *testing.T) {
	decoder := wieganddecoder.NewWiegandDecoder()

	t.Run("Valid WG26 decode", func(t *testing.T) {
		// Example of a valid WG26 bit slice (first and last bits are parity bits)
		bits := bitsFromString("10000000000101111010110011") // 26 bits
		expected := uint64(12121)                            // Example decoded value
		result, err := decoder.Decode(bits)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("Expected %d, got %d\n(%0b)\n(%0b)", expected, result, expected, result)
		}
	})

	t.Run("Valid WG34 decode", func(t *testing.T) {
		// Example of a valid WG34 bit slice (first and last bits are parity bits)
		bits := bitsFromString("1000000000010111101011001000110001") // 34 bits
		expected := uint64(3103000)                                  // Example decoded value
		result, err := decoder.Decode(bits)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("Expected %d, got %d\n(%0b)\n(%0b)", expected, result, expected, result)
		}
	})

	t.Run("Invalid bit length (unsupported)", func(t *testing.T) {
		bits := bitsFromString("101010101010") // Not 26 or 34 bits
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected error for unsupported bit length, got nil")
		}
	})

	t.Run("Parity check failed (WG26)", func(t *testing.T) {
		// Modify the bits to fail the parity check (e.g., flip a bit)
		bits := bitsFromString("10000000000101111010110010") // 26 bits, incorrect parity
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected parity check failure, got no error")
		}
	})

	t.Run("Parity check failed (WG34)", func(t *testing.T) {
		// Modify the bits to fail the parity check (e.g., flip a bit)
		bits := bitsFromString("1000000000010111101011001000110000") // 34 bits, incorrect parity
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected parity check failure, got no error")
		}
	})
}
