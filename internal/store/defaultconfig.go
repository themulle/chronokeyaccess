package store

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
			Name:       "daily",
			CronString: "0 0 0 * * * *",
			Duration:   24*time.Hour - time.Nanosecond,
			OneTimePin: true,
		},
		{
			Name:       "business hours",
			CronString: "0 0 7,11,15,19 * * * *",
			Duration:   5 * time.Hour,
			OneTimePin: true,
		},
	}

	return store
}
