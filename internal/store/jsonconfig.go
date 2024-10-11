package store

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadOrInitJsonConfiguration[T any](jsonFileName string, defaultConfig T) (T, error) {
	out := new(T)
	data, err := os.ReadFile(jsonFileName)
	if os.IsNotExist(err) {
		data, err = json.MarshalIndent(defaultConfig, "", "   ")
		if err == nil {
			err = os.WriteFile(jsonFileName, data, 0700)
		}
		if err != nil {
			err = fmt.Errorf("error generating default config: %w", err)
		}
	}
	if err == nil {
		err = json.Unmarshal(data, out)
	}

	return *out, err
}
