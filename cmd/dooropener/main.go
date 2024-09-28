package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/themulle/chronokeyaccess/pkg/dooropener"
)

func main() {
	relayPin := flag.Int("relaypin", 18, "relay pin number")
	closedState := flag.Bool("closedstate", true, "default closed state")
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, error)")
	open := flag.Bool("open", true, "initialize opener to default state")
	// Parse command-line flags
	flag.Parse()

	setupLogger(*logLevelFlag)

	slog.Info("starting door opener")
	opener := dooropener.NewDoorOpener(*relayPin, *closedState)
	opener.InitAsOutput()
	if *open {
		opener.OpenDoor()
	}

	slog.Info("done")
}

func startShutdownListener(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("shutdown")
		cancel()
	}()
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
