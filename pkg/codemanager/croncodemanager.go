package codemanager

import (
	"fmt"
	"time"
)

type cronCodeManager struct {
	codeManagerBase
	Slots CronCodeSlots
}

func NewCronCodeManager(password string, slots CronCodeSlots) *cronCodeManager {
	retval := &cronCodeManager{
		codeManagerBase: codeManagerBase{Password: password},
		Slots:           slots,
	}

	if err := retval.Init(); err != nil {
		fmt.Println(err.Error())
	}

	return retval
}

func (ecm *cronCodeManager) Init() error {
	return ecm.Slots.Init()
}

func (ecm *cronCodeManager) GetEntranceCodes(dayTime time.Time) EntranceCodes {
	retval := EntranceCodes{}

	for _, slot := range ecm.Slots {
		startTime := ecm.truncateToDay(dayTime).Add(time.Hour * -24)
		endTime := ecm.truncateToDay(dayTime).Add(time.Hour * 24)
		dayStart := ecm.truncateToDay(dayTime)

		pinCode := uint(0)
		if slot.Type == OneTimePin || slot.Type == SeriesPin {
			pinCode = ecm.CalculatePinCode(slot.CronString)
		} else if slot.Type == PersonalPin && slot.PinCode > 0 {
			pinCode = slot.PinCode
		}

		nextTime := slot.cronExpression.Next(startTime)
		for ; nextTime.Before(endTime) && nextTime.After(startTime); nextTime = slot.cronExpression.Next(nextTime) {
			nextEndTime := nextTime.Add(slot.Duration)

			if nextEndTime.After(dayStart) && nextTime.Before(dayStart.Add(time.Hour*24)) {
				if slot.Type == OneTimePin {
					pinCode = ecm.CalculatePinCode(nextTime.Format("2006-01-02 15:04"))
				}

				retval = append(retval, EntranceCode{
					Start:   nextTime,
					Stop:    nextEndTime,
					PinCode: pinCode,
					Slot:    slot,
				})
			}
		}
	}
	return retval
}

func (ecm *cronCodeManager) IsValid(currentTime time.Time, pinCode uint) bool {
	codes := ecm.GetEntranceCodes(currentTime)
	for _, c := range codes {
		if c.IsInside(currentTime) && c.PinCode == pinCode {
			return true
		}
	}
	return false
}
