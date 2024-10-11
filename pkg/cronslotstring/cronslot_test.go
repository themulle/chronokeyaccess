package cronslotstring_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/themulle/chronokeyaccess/pkg/cronslotstring"
)

func TestConvertToCronSlot(t *testing.T) {
	testCases := []struct {
		input          string
		expectedCron   string
		expectedDur    time.Duration
		expectedErrMsg string
	}{
		{"Mo 00:00-23:59", "0 00 00 * * 1 *", 24*time.Hour-1, ""},
		{"14:00-16:00 Mo,Di,Mi Okt-Mar", "0 00 14 * 10-3 1,2,3 *", 120 * time.Minute, ""},
		{"12:30-13:00 Sa Mar,Apr,Mai", "0 30 12 * 3,4,5 6 *", 30 * time.Minute, ""},
		{"14:00-16:00 Mo,Di,Mi Okt-Mar 2024,2025", "0 00 14 * 10-3 1,2,3 2024,2025", 120 * time.Minute, ""},
		{"Okt-Mar 14:00-16:00 Mo,Di,Mi", "0 00 14 * 10-3 1,2,3 *", 120 * time.Minute, ""},
		{"Okt - Mar Mo,Di,Mi 14:00 - 16:00", "0 00 14 * 10-3 1,2,3 *", 120 * time.Minute, ""},

		{"14:00 - 16:00", "0 00 14 * * * *", 120 * time.Minute, ""},
		{"14:00 - 16:00 2025", "0 00 14 * * * 2025", 120 * time.Minute, ""},
		{"Mo", "0 00 00 * * 1 *", 24*time.Hour-1, ""},
		{"fdasfdsa", "0 00 00 * * * *", 24*time.Hour-1, "unrecognized cron expression"},
		
	}

	for _, tc := range testCases {
		cron, dur, err := cronslotstring.ParseCronSlotString(tc.input)

		if err != nil && tc.expectedErrMsg == "" {
			fmt.Printf("unexpected error: %v\n", err)
			t.Fail()
		}
		if cron != tc.expectedCron {
			fmt.Printf("cron mismatch: expected %s, got %s\n", tc.expectedCron, cron)
			t.Fail()
		}
		if dur != tc.expectedDur {
			fmt.Printf("duration mismatch: expected %v, got %v\n", tc.expectedDur, dur)
			t.Fail()
		}
	}
}

func TestReverseCronSlot(t *testing.T) {
	tests := []struct {
		cronStr  string
		duration time.Duration
		expected string
	}{
		{
			cronStr:  "0 00 00 * * * *",
			duration: time.Hour*24-1,
			expected: "*",
		},
		{
			cronStr:  "0 30 08 * 1 1 *",
			duration: time.Minute*60,
			expected: "08:30-09:30 Mo Jan",
		},
		{
			cronStr:  "0 15 14 * 5 5 2024",
			duration: time.Minute*150,
			expected: "14:15-16:45 Fr Mai 2024",
		},
		{
			cronStr:  "0 15 14 * * 1-5 2024",
			duration: time.Minute*150,
			expected: "14:15-16:45 Mo-Fr 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.cronStr, func(t *testing.T) {
			result, err := cronslotstring.GenerateCronSlotString(tt.cronStr, tt.duration)
			if result!=tt.expected {
				fmt.Println(result, tt.expected)
				t.Fail()
			}
			if err!=nil {
				fmt.Println(err)
				t.Fail()
			}
			
			
		})
	}
}