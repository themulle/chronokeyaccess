package codemanager

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

type CodeManager interface {
	GetEntranceCodes(time.Time) EntranceCodes     //get all entrance codes of this day
	IsValid(time.Time, uint) (bool, EntranceCode) //check if the code is valid at this time
}

type EntranceSlot interface {
	GetName() string
	GetDescription() string
	GetType() CronCodeSlotType
	GetPinCode() uint
}

type codeManagerBase struct {
	Password string
}

func (codeManagerBase) truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (ecm *codeManagerBase) CalculatePinCode(timeString string) uint {
	hash := sha256.Sum256(append([]byte(ecm.Password), []byte(timeString)...))
	pinCode := binary.BigEndian.Uint64(hash[:8]) % 10000
	return uint(pinCode)
}
