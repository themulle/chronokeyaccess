package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/themulle/chronokeyaccess/internal/store"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

type CodeRequest struct {
	DayTime    []string `form:"daytime"`
	CodeType   string   `form:"codetype"`
	ExactMatch bool     `form:"exactmatch"`
}

func getCodes(cr CodeRequest) (codemanager.EntranceCodes, error) {
	retval := codemanager.EntranceCodes{}
	var cm codemanager.CodeManager
	if len(cr.DayTime) == 0 {
		return retval, fmt.Errorf("no daytime in request")
	}
	cr.CodeType = strings.TrimSpace(strings.ToLower(cr.CodeType))

	if codeManagerStore, err := store.LoadConfiguration("config.json", "personalcodes.csv", true); err != nil {
		return retval, err
	} else if cm, err = codemanager.InitFromStore(codeManagerStore); err != nil {
		return retval, err
	}
	for _, dayTimeString := range cr.DayTime {
		var dayTime time.Time
		var err error
		if dayTime, err = time.Parse("2006-01-02T15:04:05-0700", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("2006-01-02", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("2006-01-02T15:04", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("01.02.2006", dayTimeString); err == nil {
		} else {
			return retval, fmt.Errorf("invalid date format: %s", dayTimeString)
		}
		log.Printf("code request for: %s", dayTime.String())
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

func getSeriesPin(cr CodeRequest) (codemanager.EntranceCodes, error) {
	retval := codemanager.EntranceCodes{}
	var cm codemanager.CodeManager

	if codeManagerStore, err := store.LoadConfiguration("config.json", "personalcodes.csv", true); err != nil {
		return retval, err
	} else if cm, err = codemanager.InitFromStore(codeManagerStore); err != nil {
		return retval, err
	} else {

		/*for _,s := range codeManagerStore.Slots {
			if ccs, err:=s.(*codemanager.CronCodeSlot); err==nil {

			}
		}*/
	}

	for _, dayTimeString := range cr.DayTime {
		var dayTime time.Time
		var err error
		if dayTime, err = time.Parse("2006-01-02T15:04:05-0700", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("2006-01-02", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("2006-01-02T15:04", dayTimeString); err == nil {
		} else if dayTime, err = time.Parse("01.02.2006", dayTimeString); err == nil {
		} else {
			return retval, fmt.Errorf("invalid date format: %s", dayTimeString)
		}
		log.Printf("code request for: %s", dayTime.String())
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
		retval = append(retval, dayCodes...)
	}

	retval.Sort()
	return retval.Uniq(), nil
}
