package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
	"github.com/themulle/chronokeyaccess/pkg/gpioout"
	"github.com/themulle/chronokeyaccess/pkg/wiegand"
	"github.com/themulle/chronokeyaccess/pkg/wiegand/wieganddecoder"
)

const (
	defaultConfigFile    = "config.json"
	defaultPinFile      = "pins.csv"
	defaultAccessLogFile = "accesslog.csv"
	defaultLogLevel     = "info"
	commandTimeout    = 5 * time.Second
	doorOpenDuration  = 5 * time.Second
	buzzerDuration    = time.Second
	ledDuration       = time.Second
)

var (
	configFileName      string
	personalPinFileName string
	accessLogFile      string
)

func setupLogger(logLevelFlag string) error {
	// Definiere eine Map für die Zuordnung von Strings zu slog.Level
	logLevels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	// Versuche, das Log-Level aus der Map zu erhalten
	logLevel, ok := logLevels[logLevelFlag]
	if !ok {
		return fmt.Errorf("ungültiges Log-Level: %s. Verwenden Sie 'debug', 'info', 'warn' oder 'error'", logLevelFlag)
	}

	// Konfiguriere den Logger mit dem gewählten Level
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Logge die erfolgreiche Konfiguration
	slog.Info("Logger erfolgreich konfiguriert", 
		"level", logLevelFlag,
		"handler", "JSONHandler")

	return nil
}

func init() {
	flag.StringVar(&configFileName, "c", defaultConfigFile, "config file name")
	flag.StringVar(&personalPinFileName, "p", defaultPinFile, "personal pin csv file name")
	flag.StringVar(&accessLogFile, "accesslog", defaultAccessLogFile, "accesslog file")
}

func main() {
	data0pin := flag.Int("data0", 14, "data0 pin number")
	data1pin := flag.Int("data1", 15, "data1 pin number")
	logLevelFlag := flag.String("loglevel", defaultLogLevel, "log level (debug, info, warn, error)")
	flag.Parse()

	if err := setupLogger(*logLevelFlag); err != nil {
		slog.Error("Fehler beim Einrichten des Loggers", "error", err)
		os.Exit(1)
	}

	slog.Debug("Konfiguration geladen", 
		"configFile", configFileName,
		"pinFile", personalPinFileName)

	reader := wiegand.NewWiegandReader(*data0pin, *data1pin)
	decoder := wieganddecoder.NewCommonDecoder()

	slog.Info("Wiegand-Leser wird gestartet", 
		"data0pin", *data0pin,
		"data1pin", *data1pin)

	if err := reader.Start(); err != nil {
		slog.Error("Fehler beim Initialisieren des Wiegand-Lesers", "error", err)
		os.Exit(1)
	}
	defer reader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startShutdownListener(cancel)

	processInput(ctx, reader, decoder)

	slog.Info("Done")
}

func processInput(ctx context.Context, reader *wiegand.WiegandReader, decoder *wieganddecoder.CommonDecoder) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Verarbeitung beendet aufgrund von Kontextabbruch")
			return
		case <-time.After(100 * time.Millisecond):
			bits, start, stop := reader.GetCache()
			if decoder.CheckInputComplete(bits, start, stop) {
				reader.ResetCache()
				
				number, err := decoder.Decode(bits)
				logInput(number, bits, start, stop, err)
				
				if err := openDoor(number); err != nil {
					slog.Error("Fehler beim Öffnen der Tür", "error", err, "code", number)
				}
			}
		}
	}
}

func logInput(number uint64, bits []bool, start, stop time.Time, err error) {
	slog.Info("Input received",
		"number", number,
		"cache", wieganddecoder.BitsToString(bits),
		"error", err,
		"start", start,
		"stop", stop)

	fmt.Printf("%d %s %s %s %t\n", number, start.Format(time.RFC3339), stop.Format(time.RFC3339), wieganddecoder.BitsToString(bits), err == nil)
}

func startShutdownListener(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		slog.Info("Shutdown signal received")
		cancel()
	}()
}


func openDoor(code uint64) error {
	slog.Info("Checking code", "code", code)
	
	result, err := verifyCode(code)
	if err != nil {
		return fmt.Errorf("error during code verification: %w", err)
	}

	if result == "ok" {
		slog.Info("Code verified, opening door...")
		
		
		//go activatePin(22, 1, buzzerDuration) // Activate buzzer (Pin 22)
		go activatePin(24, 1, ledDuration)    // Activate LED (Pin 24)
		
		if err := activatePin(23, 1, doorOpenDuration); err != nil { // Activate relay (Pin 23)
			return fmt.Errorf("error activating the relay: %w", err)
		}
	} else {
		slog.Info("Code verification failed", "result", result)
		
		if err := activatePin(22, 1, buzzerDuration); err != nil {
			return fmt.Errorf("error activating the buzzer: %w", err)
		}
	}

	return nil
}

func verifyCode(code uint64) (string, error) {
	cm, err := initCodeManager()
	if err != nil {
		return "", fmt.Errorf("error initializing the CodeManager: %w", err)
	}

	valid, details := cm.IsValid(time.Now(), uint(code))

	slog.Info("code verification", "code", code, "valid", valid)
	
	
	if valid {
		logAccess(code, details.Slot.GetName())
		return "ok", nil
	}

	logAccess(code, "invalid")
	return "invalid", nil
}

func initCodeManager() (codemanager.CodeManager, error) {
	codeManagerStore, err := store.LoadConfiguration(configFileName, personalPinFileName)
	if err != nil {
		return nil, err
	}
	return codemanager.InitFromStore(codeManagerStore)
}

func logAccess(code uint64, userString string) {
	if accessLogFile == "" {
		return
	}

	logEntry := fmt.Sprintf("%s,%d,%s\n", time.Now().Format(time.RFC3339), code, userString)
	if err := appendAccessLog(accessLogFile, logEntry); err != nil {
		slog.Error("Error writing to the access log", "error", err)
	}
}

func appendAccessLog(fileName, content string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func activatePin(pin, state int, duration time.Duration) error {
	gpio := gpioout.NewGpioOut(pin, state)
	if err := gpio.InitAsOutput(); err != nil {
		return fmt.Errorf("error initializing Pin %d: %w", pin, err)
	}
	return gpio.ActivateFor(duration)
}
