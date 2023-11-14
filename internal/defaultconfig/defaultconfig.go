package defaultconfig

import (
	"log"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func GetDefualtConfig() codemanager.CodeManagerStore {
	store := codemanager.CodeManagerStore{
		Password: "1234567890",
	}

	if password, err := password.Generate(10, 4, 2, false, true); err == nil {
		store.Password = password
	} else {
		log.Printf("error generating password: %s", err)
	}

	store.Slots = codemanager.CronCodeSlots{
		{
			Description: "daily",
			CronString:  "0 0 0 * * * *",
			Duration:    24*time.Hour - time.Nanosecond,
			UseSameCode: false,
		},
		{
			Description: "business hours",
			CronString:  "0 0 8 * * * *",
			Duration:    14 * time.Hour,
			UseSameCode: true,
		},
		{
			Description: "evening",
			CronString:  "0 0 16 * * * *",
			Duration:    6 * time.Hour,
			UseSameCode: true,
		},
	}

	return store
}
