package codemanager

type CodeManagerStore struct {
	Password string
	Slots    CronCodeSlots
}

func InitFromStore(store CodeManagerStore) (CodeManager, error) {
	cm := NewCronCodeManager(store.Password, store.Slots)
	err := cm.Init()
	return cm, err
}
