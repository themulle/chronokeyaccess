package codemanager

import "github.com/gorhill/cronexpr"

type CodeManagerStore struct {
	Password      string
	Slots         CronCodeSlots
	PersonalCodes PersonalCodes `json:"-"`
}

func InitFromStore(store CodeManagerStore) (CodeManager, error) {

	for _, pc := range store.PersonalCodes {
		for _, s := range store.Slots {
			if s.Type == PersonalPin && s.Name == pc.SlotName {
				store.Slots = append(store.Slots, CronCodeSlot{
					Name:           pc.Name,
					CronString:     s.CronString,
					Duration:       s.Duration,
					Type:           PersonalPin,
					PinCode:        pc.PinCode,
					cronExpression: &cronexpr.Expression{},
				})
			}
		}
	}

	cm := NewCronCodeManager(store.Password, store.Slots)
	err := cm.Init()
	return cm, err
}
