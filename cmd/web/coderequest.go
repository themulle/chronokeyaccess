package main

import (
	"fmt"
	"strings"

	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/accesslog"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
	"github.com/themulle/chronokeyaccess/pkg/dateparser"
)

type CodeRequest struct {
	DayTime    []string `form:"daytime"`
	CodeType   string   `form:"codetype"`
	ExactMatch bool     `form:"exactmatch"`
}

func getCodes(cr CodeRequest, cm codemanager.CodeManager ) (codemanager.EntranceCodes, error) {
	retval := codemanager.EntranceCodes{}
	if len(cr.DayTime) == 0 {
		return retval, fmt.Errorf("no daytime in request")
	}
	cr.CodeType = strings.TrimSpace(strings.ToLower(cr.CodeType))

	
	for _, dayTimeString := range cr.DayTime {
		dayTime, err := dateparser.Parse(dayTimeString)
		if err != nil {
			return retval, fmt.Errorf("invalid date format: %s", dayTimeString)
		}
		dayCodes := cm.GetEntranceCodes(dayTime)
		if cr.ExactMatch {
			tmp := codemanager.EntranceCodes{}
			for _, d := range dayCodes {
				if d.Start == dayTime {
					tmp = append(tmp, d)
				}
			}
			dayCodes = tmp
		}
		{
			//filter one-time/seriespin
			tmp := codemanager.EntranceCodes{}
			for _, d := range dayCodes {
				if d.PinCode == 0 {
					continue
				}

				if cr.CodeType == "onetimepin" && d.Slot.GetType() == codemanager.OneTimePin {
					tmp = append(tmp, d)
				} else if cr.CodeType == "seriespin" && d.Slot.GetType() == codemanager.SeriesPin {
					tmp = append(tmp, d)
				} else if cr.CodeType == "personalpin" && d.Slot.GetType() == codemanager.PersonalPin {
					tmp = append(tmp, d)
				}
			}
			dayCodes = tmp
		}

		retval = append(retval, dayCodes...)
	}

	retval.Sort()
	return retval.Uniq(), nil
}

func getAccessLogs() (accesslog.AccessLogs, error) {
	var retval accesslog.AccessLogs
	var err error

	if retval, err = store.LoadAccessLogCSV("accesslog.csv"); err != nil {
		return retval, err
	}

	retval.Sort()
	return retval, nil
}
