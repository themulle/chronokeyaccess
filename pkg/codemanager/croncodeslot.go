package codemanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/themulle/cronexpr"
)

type CronCodeSlot struct {
	Name           string
	CronString     string
	Duration       time.Duration
	Type           CronCodeSlotType
	PinCode        uint
	ValidTo			time.Time
	cronExpression *cronexpr.Expression
}

func (c CronCodeSlot) GetName() string {
	return c.Name
}


func (c CronCodeSlot) GetDescription() string {
	return c.CronString
}

func (c CronCodeSlot) GetType() CronCodeSlotType {
	return c.Type
}

func (c CronCodeSlot) GetPinCode() uint {
	return c.PinCode
}

func (ces *CronCodeSlot) Init() error {
	var durationError, expressionError error

	if ces.Duration < time.Minute {
		durationError = fmt.Errorf("duaration must be at least one minute")
	}
	if ces.Duration >= time.Hour*24 {
		durationError = fmt.Errorf("duaration must be less than 24h")
	}

	ces.cronExpression, expressionError = cronexpr.Parse(ces.CronString)
	return errors.Join(durationError, expressionError)
}

type CronCodeSlots []CronCodeSlot

func (ccs CronCodeSlots) Init() error {
	var allErrors error
	for i := range ccs {
		if err := ccs[i].Init(); err != nil {
			allErrors = errors.Join(allErrors, fmt.Errorf("error in CronCode %s: %s", ccs[i].CronString, err))
		}
	}
	return allErrors
}
