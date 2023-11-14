package codemanager

import (
	"encoding/json"
	"fmt"
)

type CodeManagerStore struct {
	Password string
	Slots    CronCodeSlots
}

func Load(data []byte) (CodeManager, error) {
	store := &CodeManagerStore{}

	if err := json.Unmarshal(data, store); err != nil {
		return nil, err
	}

	if err := store.Slots.Init(); err != nil {
		return nil, err
	}

	cm := NewCronCodeManager(store.Password, store.Slots)
	return cm, nil
}

func MarshalCodeManager(cm CodeManager) ([]byte, error) {
	ccm, successful := cm.(*cronCodeManager)
	if successful {
		return nil, fmt.Errorf("not a cronCodeManager type")
	}
	store := CodeManagerStore{
		Password: ccm.Password,
		Slots:    ccm.Slots,
	}

	return MarshalCodeManagerStore(store)
}

func MarshalCodeManagerStore(store CodeManagerStore) ([]byte, error) {
	return json.MarshalIndent(store, "", "   ")
}
