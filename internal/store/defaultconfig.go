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
			Type:       codemanager.OneTimePin,
		},
		{
			Name:       "business hours",
			CronString: "0 0 7,11,15,19 * * * *",
			Duration:   5 * time.Hour,
			Type:       codemanager.OneTimePin,
		},
	}

	store.Slots = append(store.Slots, codemanager.CronCodeSlots{
		{
			Name:       "monday evening",
			CronString: "0 0 17 * * monday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "tuesday evening",
			CronString: "0 0 17 * * tuesday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "wednesday evening",
			CronString: "0 0 17 * * wednesday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "thursday evening",
			CronString: "0 0 17 * * thursday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "friday evening",
			CronString: "0 0 17 * * friday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "saturday evening",
			CronString: "0 0 17 * * saturday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
		{
			Name:       "sunday evening",
			CronString: "0 0 17 * * sunday *",
			Duration:   7 * time.Hour,
			Type:       codemanager.SeriesPin,
		},
	}...)

	return store
}
