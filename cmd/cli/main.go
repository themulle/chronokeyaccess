package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

var (
	configFileName      string
	personalPinFileName string
	pinCode             uint
	showCurrentCodes    uint
	accessLogFile       string
)

func init() {
	flag.StringVar(&configFileName, "c", "config.json", "config file name")
	flag.StringVar(&personalPinFileName, "p", "personalcodes.csv", "personal pin csv file name")
	flag.StringVar(&accessLogFile, "accesslog", "accesslog.csv", "accesslog file")
	flag.UintVar(&pinCode, "t", 0, "test pin code")
	flag.UintVar(&showCurrentCodes, "s", 0, "show codes for the next n days")
	flag.Parse()
}

func main() {
	cm, err := func() (codemanager.CodeManager, error) {
		codeManagerStore, err := store.LoadConfiguration(configFileName, personalPinFileName, true)
		if err != nil {
			return nil, err
		}
		return codemanager.InitFromStore(codeManagerStore)
	}()

	if err != nil {
		fmt.Printf("configuration error: %s", err)
		os.Exit(1)
	}

	if showCurrentCodes > 0 {
		allCodes := codemanager.EntranceCodes{}
		for i := 0; i < int(showCurrentCodes); i++ {
			offset := 24 * time.Hour * time.Duration(i)
			allCodes = append(allCodes, cm.GetEntranceCodes(time.Now().Local().Add(offset))...)
		}
		allCodes.Sort()
		allCodes = allCodes.Uniq()
		fmt.Println(allCodes.String())
		os.Exit(0)
	}

	if pinCode > 0 {
		valid, details := cm.IsValid(time.Now(), pinCode)
		userString := "invalid"

		if valid {
			userString = details.Slot.GetName()
			fmt.Println("ok")
		} else {
			fmt.Println("invalid")
		}

		if accessLogFile != "" {
			appendAccessLog(accessLogFile, fmt.Sprintf("%s,%d,%s\n", time.Now().Format("2006-01-02T15:04:05"), pinCode, userString))
		}

		os.Exit(0)
	}

	flag.PrintDefaults()
}

func appendAccessLog(fileName, content string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return err
	}

	return nil
}
