package codemanager

import "encoding/json"

type CronCodeSlotType int

const (
	Undefined CronCodeSlotType = iota
	OneTimePin
	SeriesPin
	PersonalPin
)

func (t CronCodeSlotType) String() string {
	switch t {
	case OneTimePin:
		return "OneTimePin"
	case SeriesPin:
		return "SeriesPin"
	case PersonalPin:
		return "PersonalPin"
	default:
		return "undefined"
	}
}

func (t CronCodeSlotType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *CronCodeSlotType) UnmarshalJSON(data []byte) (err error) {
	var typeString string
	if err := json.Unmarshal(data, &typeString); err != nil {
		return err
	}
	newVal := CronCodeSlotType(Undefined)
	switch typeString {
	case "OneTimePin":
		newVal = OneTimePin
	case "SeriesPin":
		newVal = SeriesPin
	case "PersonalPin":
		newVal = PersonalPin
	default:
		newVal = Undefined
	}
	*t = newVal
	return nil
}
