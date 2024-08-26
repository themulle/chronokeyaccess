package accesslog

import (
	"sort"
	"time"
)

type AccessLog struct {
	Ts      time.Time //time pin was entered
	PinCode uint      //pin code used
	Status  string    //status (failure or user/slot name)
}

func (ec AccessLog) Equals(other AccessLog) bool {
	return ec.Ts.Equal(other.Ts) && ec.PinCode == other.PinCode
}

type AccessLogs []AccessLog

func (ec AccessLogs) Sort() {
	sort.Slice(ec, func(i, j int) bool {
		return ec[i].Ts.Before(ec[j].Ts)
	})
}
