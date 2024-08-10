package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func LoadConfiguration(configFileName string, createDefaultUnlessExistent bool) (codemanager.CodeManagerStore, error) {
	var store codemanager.CodeManagerStore
	data, err := os.ReadFile(configFileName)
	if os.IsNotExist(err) {
		data, err = json.Marshal(GetDefualtConfig())
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
	return store, err
}
