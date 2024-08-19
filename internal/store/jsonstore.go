package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func LoadConfiguration(configFileName string, personalCodeFileName string, createDefaultUnlessExistent bool) (codemanager.CodeManagerStore, error) {
	var store codemanager.CodeManagerStore
	data, err := os.ReadFile(configFileName)
	if os.IsNotExist(err) {
		data, err = json.MarshalIndent(GetDefualtConfig(), "", "   ")
		if err == nil {
			err = os.WriteFile(configFileName, data, 0700)
		}
		if err != nil {
			err = fmt.Errorf("error generating default config: %w", err)
		}
	}
	if err == nil {
		err = json.Unmarshal(data, &store)
	}

	if err != nil {
		return store, err
	}

	{
		var pc codemanager.PersonalCodes
		pc, err = LoadPersonalCodeCSV(personalCodeFileName)
		if os.IsNotExist(err) {
			if err = WritePersonalCodeCSV(GetDefualtConfig().PersonalCodes, personalCodeFileName); err != nil {
				err = fmt.Errorf("error generating personal code file: %w", err)
			}
		} else if err == nil {
			store.PersonalCodes = pc
		}
	}

	return store, err
}
