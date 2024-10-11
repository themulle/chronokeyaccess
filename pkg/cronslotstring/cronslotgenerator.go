package cronslotstring

import (
	"fmt"
	"strings"
	"time"

	"github.com/themulle/cronexpr"
)

// ReverseCronSlot takes a cron expression and reverses it to generate the original input string
func (csg* CronSlotGenerator) GenerateCronSlotString(cronStr string, duration time.Duration) (string, error) {
	// Validate cron expression
	_, err := cronexpr.Parse(cronStr)
	if err != nil {
		return "", fmt.Errorf("invalid cron string: %v", err)
	}

	parts := strings.Split(cronStr, " ")
	if len(parts) != 7 {
		return "", fmt.Errorf("invalid cron format")
	}

	minutes := parts[1]
	hours := parts[2]
	months := csg.reverseExtractMonths(parts[4])
	weekdays := csg.reverseExtractDays(parts[5])
	years := parts[6]

	// Construct time range from hours and minutes
	timeRange:=fmt.Sprintf("%s:%s", hours, minutes)
	timeRange=fmt.Sprintf("%s-%s", timeRange, calculateEnd(timeRange,duration))

	// Build the reverse input string
	
	var result string
	if timeRange!="00:00-23:59" {
		result+=timeRange+ " "
	}
	if weekdays != "*" {
		result += weekdays + " "
	}
	if months != "*" {
		result += months + " "
	}
	
	if years != "*" {
		result += years + " "
	}

	if result=="" {
		result="*"
	}
	return strings.TrimSpace(result), nil
}

func calculateEnd(startTime string, duration time.Duration) string {
	start, err := time.Parse("15:04", startTime)
	if err!=nil {
		return ""
	}

	end := start.Add(duration)

	return end.Format("15:04")
}

func (csg* CronSlotGenerator) reverseExtractDays(input string) string {
	for dayNumber, dayName := range reverseMap(csg.LocaleWeekDayMap) {
		input = strings.ReplaceAll(input, dayNumber, dayName)
	}
	return input
}

func  (csg* CronSlotGenerator) reverseExtractMonths(input string) string {
	for monthNumber, monthName := range reverseMap(csg.LocaleMonthNameMap) {
		input = strings.ReplaceAll(input, monthNumber, monthName)
	}
	return input
}

func getMapKeys(in ...map[string]string) []string {
	keys:=[]string{}

	for _, inMap := range in {
        for k := range inMap {
			keys = append(keys, k)
		}
    }
	return keys

}

func reverseMap(in map[string]string) map[string]string {
	reversed := map[string]string{}
	for k,v := range in {
		reversed[v]=k
	}
	return reversed
}
