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
)

func init() {
	flag.StringVar(&configFileName, "c", "config.json", "config file name")
	flag.StringVar(&personalPinFileName, "p", "personalpin.csv", "personal pin csv file name")
	flag.UintVar(&pinCode, "t", 0, "test pin code")
	flag.UintVar(&showCurrentCodes, "s", 1, "show codes for the next n days")
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

	if pinCode > 0 {
		if cm.IsValid(time.Now(), pinCode) {
			fmt.Println("ok")
			os.Exit(0)
		} else {
			fmt.Println("invalid")
			os.Exit(0)
		}
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

	flag.PrintDefaults()
}
