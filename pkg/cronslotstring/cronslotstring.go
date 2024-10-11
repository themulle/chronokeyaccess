package cronslotstring

import (
	"fmt"
	"time"

	"github.com/goodsign/monday"
)

type CronSlotGenerator struct {
	LocaleMonthNameMap map[string]string
	LocaleWeekDayMap map[string]string

	MonthNameMap map[string]string
	WeekDayMap map[string]string
}

func NewCronSlotGenerator() *CronSlotGenerator {
	csg := &CronSlotGenerator{
		MonthNameMap:       map[string]string{},
		LocaleMonthNameMap: map[string]string{},
		WeekDayMap:         map[string]string{},
		LocaleWeekDayMap:   map[string]string{},
	}
	csg.addWeekDaysToMap(csg.WeekDayMap,monday.LocaleEnUS)
	csg.addWeekDaysToMap(csg.LocaleWeekDayMap,monday.LocaleEnUS)
	csg.addMonthNamesToMap(csg.MonthNameMap,monday.LocaleEnUS)
	csg.addMonthNamesToMap(csg.LocaleMonthNameMap,monday.LocaleEnUS)

	return csg
}

func (c * CronSlotGenerator) SetLocale(l monday.Locale) {
	c.LocaleWeekDayMap=map[string]string{}
	c.addWeekDaysToMap(c.LocaleWeekDayMap,l)
	c.LocaleMonthNameMap=map[string]string{}
	c.addMonthNamesToMap(c.LocaleMonthNameMap,l)
}

func (c * CronSlotGenerator) addWeekDaysToMap(target map[string]string,  l monday.Locale) { 
	for i:=0; i<7; i++ {
		tDay:=time.Date(0,1,1+i,0,0,0,0,time.UTC)
		tDay=tDay.Add(time.Hour*24)
		translated := monday.Format(tDay, "Mon", l)
		target[translated]=fmt.Sprintf("%d",i)
	}
}


func (c * CronSlotGenerator) addMonthNamesToMap(target map[string]string,  l monday.Locale) { 
	for i:=time.January; i<=time.December; i++ {
		tDay:=time.Date(0,i,1,0,0,0,0,time.UTC)
		translated := monday.Format(tDay, "Jan", l)
		target[translated]=fmt.Sprintf("%d",i)
	}
}

