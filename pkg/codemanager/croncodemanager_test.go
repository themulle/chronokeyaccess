package codemanager_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/goodsign/monday"
	"github.com/themulle/chronokeyaccess/pkg/codemanager"
)

func TestCalculateCronPinCode_Static(t *testing.T) {
	ecm := codemanager.NewCronCodeManager("1234", codemanager.CronCodeSlots{
		{
			CronString: "0 0 16 ? 1 mon,tue,wed,thu,fri 2023",
			Duration:   time.Hour * 3,
			OneTimePin: false,
		},
		{
			CronString: "0 0 20 ? 1 1 2023",
			Duration:   time.Hour * 3,
			OneTimePin: false,
		},
		{
			CronString: "0 0 8 ? * 6 2023",
			Duration:   time.Hour * 4,
			OneTimePin: true,
		},
	})

	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	allCodes := codemanager.EntranceCodes{}
	for t1.Before(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)) {
		dayCodes := ecm.GetEntranceCodes(t1)
		allCodes = append(allCodes, dayCodes...)
		t1 = t1.Add(time.Hour * 24)
	}

	allCodes.Sort()

	if len(allCodes) != 35 {
		t.Fatalf("35 EntranceCodes expected, got %d", len(allCodes))
	}

	result := ""
	for _, c := range allCodes {
		result += fmt.Sprintf("%s - %s PIN:%04d\n", monday.Format(c.Start, "Mon 02.01.2006 15:04", monday.LocaleDeDE), c.Stop.Format("15:04"), c.PinCode)
	}
	fmt.Println(result)

	if result != "Mo 02.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mo 02.01.2023 20:00 - 23:00 PIN:6467\n"+
		"Di 03.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mi 04.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Do 05.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Fr 06.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Sa 07.01.2023 08:00 - 12:00 PIN:3494\n"+
		"Mo 09.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mo 09.01.2023 20:00 - 23:00 PIN:6467\n"+
		"Di 10.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mi 11.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Do 12.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Fr 13.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Sa 14.01.2023 08:00 - 12:00 PIN:4168\n"+
		"Mo 16.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mo 16.01.2023 20:00 - 23:00 PIN:6467\n"+
		"Di 17.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mi 18.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Do 19.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Fr 20.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Sa 21.01.2023 08:00 - 12:00 PIN:9596\n"+
		"Mo 23.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mo 23.01.2023 20:00 - 23:00 PIN:6467\n"+
		"Di 24.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mi 25.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Do 26.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Fr 27.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Sa 28.01.2023 08:00 - 12:00 PIN:4399\n"+
		"Mo 30.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Mo 30.01.2023 20:00 - 23:00 PIN:6467\n"+
		"Di 31.01.2023 16:00 - 19:00 PIN:3630\n"+
		"Sa 04.02.2023 08:00 - 12:00 PIN:0642\n"+
		"Sa 11.02.2023 08:00 - 12:00 PIN:0976\n"+
		"Sa 18.02.2023 08:00 - 12:00 PIN:4784\n"+
		"Sa 25.02.2023 08:00 - 12:00 PIN:8331\n" {
		t.Error("pincodes or times incorrect")
	}
}

func TestCalculateDynamicCronPinCode_Dynamic(t *testing.T) {
	ecm := codemanager.NewCronCodeManager("1234", codemanager.CronCodeSlots{
		{
			CronString: "0 0 8 ? * 6 2023",
			Duration:   time.Hour * 4,
			OneTimePin: true,
		},
	})

	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	allCodes := codemanager.EntranceCodes{}
	for t1.Before(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)) {
		dayCodes := ecm.GetEntranceCodes(t1)
		allCodes = append(allCodes, dayCodes...)
		t1 = t1.Add(time.Hour * 24)
	}

	allCodes.Sort()

	result := ""
	for _, c := range allCodes {
		result += fmt.Sprintf("%s - %s PIN:%04d\n", monday.Format(c.Start, "Mon 02.01.2006 15:04", monday.LocaleDeDE), c.Stop.Format("15:04"), c.PinCode)
	}
	fmt.Println(result)

	if result != "Sa 07.01.2023 08:00 - 12:00 PIN:3494\n"+
		"Sa 14.01.2023 08:00 - 12:00 PIN:4168\n"+
		"Sa 21.01.2023 08:00 - 12:00 PIN:9596\n"+
		"Sa 28.01.2023 08:00 - 12:00 PIN:4399\n"+
		"Sa 04.02.2023 08:00 - 12:00 PIN:0642\n"+
		"Sa 11.02.2023 08:00 - 12:00 PIN:0976\n"+
		"Sa 18.02.2023 08:00 - 12:00 PIN:4784\n"+
		"Sa 25.02.2023 08:00 - 12:00 PIN:8331\n" {
		t.Error("pincodes or times incorrect")
	}
}

func TestCalculateDynamicCronPinCode_Overflow(t *testing.T) {
	ecm := codemanager.NewCronCodeManager("1234", codemanager.CronCodeSlots{
		{
			CronString: "0 0 22 1 2 * 2023",
			Duration:   time.Hour * 4,
			OneTimePin: true,
		},
	})

	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	allCodes := codemanager.EntranceCodes{}
	for t1.Before(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)) {
		dayCodes := ecm.GetEntranceCodes(t1)
		allCodes = append(allCodes, dayCodes...)
		t1 = t1.Add(time.Hour * 24)
	}

	allCodes.Sort()
	allCodes = allCodes.Uniq()

	result := ""
	for _, c := range allCodes {
		result += fmt.Sprintf("%s - %s PIN:%04d\n", monday.Format(c.Start, "Mon 02.01.2006 15:04", monday.LocaleDeDE), c.Stop.Format("15:04"), c.PinCode)
	}
	fmt.Println(result)

	if result != "Mi 01.02.2023 22:00 - 02:00 PIN:6943\n" {
		t.Error("pincodes or times incorrect")
	}
}
