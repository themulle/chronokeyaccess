package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/wiegand"
	"github.com/themulle/chronokeyaccess/pkg/wiegand/wieganddecoder"
)

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

func main() {
	data0pin := flag.Int("data0pin", 14, "data0 pin number")
	data1pin := flag.Int("data1pin", 15, "data1 pin number")
	logLevelFlag := flag.String("loglevel", "info", "Log level (debug, info, error)")
	exec := flag.String("exec", "", "application to execute after pin entry")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if help != nil && *help {
		flag.PrintDefaults()
		return
	}

	setupLogger(*logLevelFlag)

	reader := wiegand.NewWiegandReader(*data0pin, *data1pin)
	decoder := wieganddecoder.NewCommonDecoder()

	slog.Info("starting wiegand reader")
	if err := reader.Start(); err != nil {
		slog.Error("Error initializing Wiegand reader", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	startShutdownListener(cancel)

	go func() {
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

				fmt.Printf("%d %s %s %s %t\n", number, start.Format(time.RFC3339), stop.Format(time.RFC3339), wieganddecoder.BitsToString(bits), err == nil)
				if exec != nil && len(strings.TrimSpace(*exec)) > 0 {
					if err := runExternalApp(*exec, fmt.Sprintf("%d", number)); err != nil {
						slog.Error("error executing application", "error", err)
					} else {
						slog.Info("executed", "app", *exec, "code", fmt.Sprintf("%d", number))
					}
				}
			}

		}
	}()

	<-ctx.Done()
	reader.Close()

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

func runExternalApp(appPath string, argument string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, appPath, argument)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out")
	}

	slog.Debug("executing application", "appPath", appPath, "argument", argument, "error", err, "stdout", string(stdout.Bytes()), "stderr", string(stderr.Bytes()))

	if err != nil {
		return fmt.Errorf("failed to run %s: %w", appPath, err)
	}

	return nil
}
