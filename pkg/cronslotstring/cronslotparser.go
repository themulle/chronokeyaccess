package cronslotstring

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/themulle/cronexpr"
)

// ConvertToCron takes an input string with a time range, days, and months, and converts it to a valid cron expression
func (csg * CronSlotGenerator) ParseCronSlotString(input string) (string, time.Duration, error) {
	// Normalize input by replacing extra spaces, and common separators
	input = normalizeInput(input)

	if input=="" || input=="*" {
		return "0 00 00 * * * *", time.Hour*24-1, nil
	}

	parts := strings.Split(input," ")
	var err error

	var hours string = "00"
	var minutes string = "00"
	var months string = "*"
	var years string = "*"
	var weekdays string = "*"

	var duration time.Duration = 24 * time.Hour - 1

	timeRegex := regexp.MustCompile(`^(\d{1,2}):(\d{2})-(\d{1,2}:\d{2})$`)
	monthRegex := regexp.MustCompile(`^(`+strings.Join(getMapKeys(csg.LocaleMonthNameMap,csg.MonthNameMap),"|")+`)`)
	weekDayRegex := regexp.MustCompile(`^(`+strings.Join(getMapKeys(csg.LocaleWeekDayMap,csg.WeekDayMap),"|")+`)`)
	yearRegex := regexp.MustCompile(`^(\d{4})`)
	for _, part := range parts {
		timeMatches := timeRegex.FindStringSubmatch(part)
		if len(timeMatches) == 4 {
			startTime := timeMatches[1]+":"+timeMatches[2]
			endTime := timeMatches[3]
			calculatedDuration, err := calculateDuration(startTime, endTime)
			if startTime=="00:00" && endTime=="23:59" {
				calculatedDuration=calculatedDuration+time.Minute-1
			}
			if err==nil {
				duration=calculatedDuration
				hours=timeMatches[1]
				minutes=timeMatches[2]
			}
		} else if monthRegex.MatchString(part) {
			months=csg.extractMonths(part)
		} else if weekDayRegex.MatchString(part) {
			weekdays=csg.extractWeekDays(part)
		} else if yearRegex.MatchString(part){
			years=extractYears(part)
		} else {
			err = fmt.Errorf("unrecognized cron string")
		}
	}

	cronString := fmt.Sprintf("0 %s %s * %s %s %s", minutes, hours, months, weekdays, years)
	if err==nil {
		_,err = cronexpr.Parse(cronString)
	}
	
	return cronString, duration, err
}

// normalizeInput prepares the input string by removing extra spaces, replacing common separators
func normalizeInput(input string) string {
	input = strings.ReplaceAll(input, " - ", "-")
	input = strings.ReplaceAll(input, " : ", ":")
	input = strings.ReplaceAll(input, " , ", ",")
	return strings.TrimSpace(input)
}

// calculateDuration calculates the time difference between start and end time in minutes
func calculateDuration(startTime, endTime string) (time.Duration, error) {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return 0, fmt.Errorf("invalid start time")
	}

	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return 0, fmt.Errorf("invalid end time")
	}

	if end.Before(start) {
		return end.Sub(start) + 24*time.Hour, nil // Time range crosses midnight
	}
	return end.Sub(start), nil
}

func (csg * CronSlotGenerator) extractWeekDays(input string) (string) {
	for dayName, dayNumber := range csg.LocaleWeekDayMap {
			input=strings.ReplaceAll(input,dayName,dayNumber)
	}
	return input
}

func (csg * CronSlotGenerator) extractMonths(input string) (string) {
	for monthName, monthNumber := range csg.LocaleMonthNameMap {
			input=strings.ReplaceAll(input,monthName,monthNumber)
	}
	return input
}

func extractYears(input string) (string) {
	return input
}
