package codemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type CronCodeSlotType int

const (
	Undefined CronCodeSlotType = iota
	OneTimePin
	SeriesPin
	PersonalPin
)

func (t CronCodeSlotType) String() string {
	switch t {
	case OneTimePin:
		return "OneTimePin"
	case SeriesPin:
		return "SeriesPin"
	case PersonalPin:
		return "PersonalPin"
	default:
		return "undefined"
	}
}

func (t CronCodeSlotType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *CronCodeSlotType) UnmarshalJSON(data []byte) (err error) {
	var typeString string
	if err := json.Unmarshal(data, &typeString); err != nil {
		return err
	}
	newVal := CronCodeSlotType(Undefined)
	switch typeString {
	case "OneTimePin":
		newVal = OneTimePin
	case "SeriesPin":
		newVal = SeriesPin
	case "PersonalPin":
		newVal = PersonalPin
	default:
		newVal = Undefined
	}
	*t = newVal
	return nil
}

type CronCodeSlot struct {
	Name           string
	CronString     string
	Duration       time.Duration
	Type           CronCodeSlotType
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
