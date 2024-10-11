package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/dooropener"
)

func main() {
	ledPin := flag.Int("ledPin", 22, "led pin number")
	relayPin := flag.Int("relaypin", 23, "relay pin number")
	buzzerPin := flag.Int("buzzerpin", 24, "buzzer pin number")
	closedState := flag.Bool("closedstate", true, "default closed state")
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, error)")
	open := flag.Bool("open", true, "initialize opener to default state")
	flag.Parse()

	setupLogger(*logLevelFlag)

	slog.Info("starting door opener")
	relay := dooropener.NewDoorOpener(*relayPin, *closedState)
	led := dooropener.NewGpioOut(*ledPin,*closedState)
	buzzer := dooropener.NewGpioOut(*buzzerPin, *closedState)
	relay.InitAsOutput()
	led.InitAsOutput()
	buzzer.InitAsOutput()
	if *open {
		go buzzer.ActivateFor(time.Millisecond*500)
		go led.ActivateFor(time.Second*2)
		relay.ActivateFor(time.Second*3)	
	}

	slog.Info("done")
}

func setupLogger(logLevelFlag string) error {
	logLevel := slog.LevelInfo
	switch logLevelFlag {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	default:
		return fmt.Errorf("Invalid log level: %s. Use 'debug', 'info', or 'error'", logLevelFlag)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	return nil
}
