package codemanager

type PersonalCode struct {
	Name         string 
	PinCode      uint
	SlotName     string
	CronCodeSlot CronCodeSlot
}

type PersonalCodes []PersonalCode
