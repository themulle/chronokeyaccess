package codemanager

import (
	"fmt"
	"sort"
	"time"

	"github.com/goodsign/monday"
)

type EntranceCode struct {
	Description string
	Start       time.Time
	Stop        time.Time
	PinCode     uint
}

func (ec EntranceCode) Equals(other EntranceCode) bool {
	return ec.Start.Equal(other.Start) && ec.Stop.Equal(other.Stop) && ec.PinCode == other.PinCode
}

func (ec EntranceCode) IsInside(t time.Time) bool {
	return ec.Start.Equal(t) || ec.Stop.Equal(t) || ec.Start.Before(t) && ec.Stop.After(t)
}

func (ec EntranceCode) String() string {
	return fmt.Sprintf("%s - %s PIN:%04d (%s)", monday.Format(ec.Start, "Mon 02.01.2006 15:04", monday.LocaleDeDE), ec.Stop.Format("15:04"), ec.PinCode, ec.Description)
}

type EntranceCodes []EntranceCode

func (ec EntranceCodes) Sort() {
	sort.Slice(ec, func(i, j int) bool {
		return ec[i].Start.Before(ec[j].Start)
	})
}

func (ec EntranceCodes) Uniq() EntranceCodes {
	uniqueSlice := EntranceCodes{}
	for i := 0; i < len(ec); i++ {
		if i == 0 || !ec[i].Equals(ec[i-1]) {
			uniqueSlice = append(uniqueSlice, ec[i])
		}
	}
	return uniqueSlice
}

func (ec EntranceCodes) String() string {
	retval := ""
	for _, c := range ec {
		retval += c.String() + "\n"
	}
	return retval
}
