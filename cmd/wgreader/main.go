package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/wiegand"
	"github.com/themulle/chronokeyaccess/pkg/wiegand/wieganddecoder"
)

func main() {
	reader := wiegand.NewWiegandReader(14, 15)
	decoder := wieganddecoder.NewCommonDecoder()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("starting wiegand reader")
	if err := reader.Start(); err != nil {
		slog.Error("Error initializing Wiegand reader", "error", err)
		return
	}

	for {
		time.Sleep(time.Millisecond * 100)
		bits, start, stop := reader.GetCache()
		if decoder.CheckInputComplete(bits, start, stop) {
			reader.ResetCache()
			number, err := decoder.Decode(bits)
			slog.Info("got input", "number", number,
				"cache", wieganddecoder.BitsToString(bits),
				"error", err,
				"start", start,
				"stop", stop)
		}

	}

	reader.Close()

	fmt.Println("quitting")

}
