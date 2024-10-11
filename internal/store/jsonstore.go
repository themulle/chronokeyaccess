package store

import (
	"fmt"
	"os"

	"github.com/goodsign/monday"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func LoadConfiguration(configFileName string, personalCodeFileName string) (codemanager.CodeManagerStore, error) {
	var store codemanager.CodeManagerStore
	var err error

	store, err=LoadOrInitJsonConfiguration(configFileName, GetDefualtConfig())
	if err != nil {
		return store, err
	}

	{
		var pc codemanager.PersonalCodes
		pc, err = LoadPersonalCodeCSV(personalCodeFileName, monday.LocaleDeDE)
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
