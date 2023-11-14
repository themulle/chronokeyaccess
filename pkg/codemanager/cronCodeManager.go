package codemanager

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type cronCodeManager struct {
	codeManagerBase
	Slots CronCodeSlots
}

func NewCronCodeManager(password string, slots CronCodeSlots) CodeManager {
	retval := &cronCodeManager{
		codeManagerBase: codeManagerBase{Password: password},
		Slots:           slots,
	}

	if err := slots.Init(); err != nil {
		fmt.Println(err.Error())
	}

	return retval
}

func (ecm *cronCodeManager) GetEntranceCodes(dayTime time.Time) EntranceCodes {
	retval := EntranceCodes{}

	for _, slot := range ecm.Slots {
		startTime := ecm.truncateToDay(dayTime).Add(time.Hour * -24)
		endTime := ecm.truncateToDay(dayTime).Add(time.Hour * 24)
		dayStart := ecm.truncateToDay(dayTime)

		pinCode := ecm.CalculatePinCode(slot.CronString)

		exp, err := cronexpr.Parse(slot.CronString)
		fmt.Printf("%+v %s\n", exp, err)

		nextTime := slot.cronExpression.Next(startTime)
		for ; nextTime.Before(endTime) && nextTime.After(startTime); nextTime = slot.cronExpression.Next(nextTime) {
			nextEndTime := nextTime.Add(slot.Duration)

			if nextEndTime.After(dayStart) && nextTime.Before(dayStart.Add(time.Hour*24)) {
				if !slot.UseSameCode {
					pinCode = ecm.CalculatePinCode(nextTime.Format("2006-01-02 15:04"))
				}

				retval = append(retval, EntranceCode{
					Start:       nextTime,
					Stop:        nextEndTime,
					PinCode:     pinCode,
					Description: slot.Description,
				})
			}
		}
	}
	return retval
}

func (ecm *cronCodeManager) IsValid(currentTime time.Time, pinCode uint) bool {
	codes := ecm.GetEntranceCodes(currentTime)
	for _, c := range codes {
		if c.IsInside(currentTime) {
			return true
		}
	}
	return false
}
