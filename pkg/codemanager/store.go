package codemanager

import "github.com/themulle/cronexpr"

type CodeManagerStore struct {
	Password      string
	Slots         CronCodeSlots
	PersonalCodes PersonalCodes `json:"-"`
}

func InitFromStore(store CodeManagerStore) (CodeManager, error) {

	for _, pc := range store.PersonalCodes {
		store.Slots=append(store.Slots, 
			CronCodeSlot{
				Name:           pc.Name,
				CronString:     pc.CronString,
				Duration:       pc.Duration,
				Type:           PersonalPin,
				PinCode:        pc.PinCode,
				cronExpression: &cronexpr.Expression{},
			},
		)
	}

	cm := NewCronCodeManager(store.Password, store.Slots)
	err := cm.Init()
	return cm, err
}
