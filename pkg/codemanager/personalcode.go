package codemanager

import "time"

type PersonalCode struct {
	Name     string
	PinCode  uint
	CronString string
	Duration time.Duration
}

type PersonalCodes []PersonalCode
