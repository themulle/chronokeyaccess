package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/themulle/chronokeyaccess/internal/defaultconfig"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

var (
	configFileName   string
	pinCode          uint
	showCurrentCodes uint
)

func init() {
	flag.StringVar(&configFileName, "c", "config.json", "config file name")
	flag.UintVar(&pinCode, "t", 0, "test pin code")
	flag.UintVar(&showCurrentCodes, "s", 300, "show codes for the next n days")
	flag.Parse()
}

func main() {
	var cm codemanager.CodeManager
	{
		data, err := os.ReadFile(configFileName)
		if os.IsNotExist(err) {
			fmt.Printf("config file %s does not exist, creating default config\n", configFileName)
			data, err = codemanager.MarshalCodeManagerStore(defaultconfig.GetDefualtConfig())
			if err != nil {
				fmt.Printf("Error: %s", err)
				os.Exit(1)
			}
			err = os.WriteFile(configFileName, data, 0700)
			if err != nil {
				fmt.Printf("Error: %s", err)
				os.Exit(1)
			}
		} else if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
		cm, err = codemanager.Load(data)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

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
