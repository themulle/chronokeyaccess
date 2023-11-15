package codemanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type CronCodeSlot struct {
	Description    string
	CronString     string
	Duration       time.Duration
	OneTimePin     bool
	cronExpression *cronexpr.Expression
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
			allErrors = errors.Join(allErrors, fmt.Errorf("error in CronCode %s: %w\n", ccs[i].CronString, err))

		}
	}
	return allErrors
}
