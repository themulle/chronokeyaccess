package wieganddecoder_test

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/wiegand/wieganddecoder"
)

// Helper function to create bit slices from a binary string (for readability)
func bitsFromString(s string) []bool {
	bits := make([]bool, len(s))
	for i, c := range s {
		if c == '1' {
			bits[i] = true
		} else {
			bits[i] = false
		}
	}
	return bits
}

func TestKeypadDecoder_CheckInputComplete(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	decoder := wieganddecoder.NewKeypadDecoder()

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

	t.Run("Timeout exceeded", func(t *testing.T) {
		start := time.Now().Add(-time.Second * 5)
		stop := time.Now().Add(-time.Second * 5)
		bits := bitsFromString("110010")
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (timeout exceeded), got %v", result)
		}
	})

	t.Run("Valid keypad input", func(t *testing.T) {
		start := time.Now().Add(-decoder.MaxTimeout - decoder.MaxTimeout)
		stop := time.Now().Add(-decoder.MaxTimeout - 100*time.Millisecond)
		bits := bitsFromString("00011001") // Expected valid BCD (9)
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (valid keypad input), got %v", result)
		}
	})

	t.Run("Short input time", func(t *testing.T) {
		start := time.Now().Add(-time.Millisecond * 500)
		stop := time.Now()
		bits := bitsFromString("110010")
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != false {
			t.Errorf("Expected false (input too short), got %v", result)
		}
	})

	t.Run("Magic", func(t *testing.T) {
		start := time.Now().Add(-decoder.MaxTimeout)
		stop := time.Now()
		bits := bitsFromString("01111010")
		result := decoder.CheckInputComplete(bits, start, stop)
		if result != true {
			t.Errorf("Expected true (magic number at end of input), got %v", result)
		}
	})
}

func TestKeypadDecoder_Decode(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	decoder := wieganddecoder.NewKeypadDecoder()

	t.Run("Valid BCD decode", func(t *testing.T) {
		bits := bitsFromString("00010010") // BCD representing "12"
		result, err := decoder.Decode(bits)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != 12 {
			t.Errorf("Expected 12, got %d", result)
		}
	})

	t.Run("Invalid BCD decode (digit > 9)", func(t *testing.T) {
		bits := bitsFromString("11110000") // Invalid BCD, digit > 9
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected error for invalid BCD digit, got nil")
		}
	})

	t.Run("Empty bit slice", func(t *testing.T) {
		bits := []bool{}
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected error for empty bit slice, got nil")
		}
	})

	t.Run("Non-multiple of 4 bits", func(t *testing.T) {
		bits := bitsFromString("101") // Not a multiple of 4
		_, err := decoder.Decode(bits)
		if err == nil {
			t.Fatal("Expected error for non-multiple of 4 bits, got nil")
		}
	})

	t.Run("MgicEnd", func(t *testing.T) {
		bits := bitsFromString("0111011101111010") // Not a multiple of 4
		result, err := decoder.Decode(bits)
		if err != nil {
			t.Fatal("Expected no error for magic number at the end")
		}
		if result != 777 {
			t.Errorf("Expected 7, got %d", result)
		}
	})
}
