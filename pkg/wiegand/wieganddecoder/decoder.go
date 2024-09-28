package wieganddecoder

import (
	"errors"
	"fmt"
	"time"
)

type Decoder interface {
	Decode(bits []bool) (uint64, error)
	CheckInputComplete(bits []bool, start time.Time, stop time.Time) bool
}

type CommonDecoder struct {
	keypadDecoder  Decoder
	wiegandDecoder Decoder
}

func NewCommonDecoder() *CommonDecoder {
	decoder := &CommonDecoder{
		keypadDecoder:  NewKeypadDecoder(),
		wiegandDecoder: NewWiegandDecoder(),
	}
	return decoder
}

func (c *CommonDecoder) Decode(bits []bool) (uint64, error) {
	var code uint64
	var err1, err2 error
	if c.wiegandDecoder != nil {
		code, err1 = c.wiegandDecoder.Decode(bits)
		if err1 == nil {
			return code, nil
		} else {
			err1 = fmt.Errorf("wiegandDecoder error: %w", err1)
		}
	}
	if c.keypadDecoder != nil {
		code, err2 = c.keypadDecoder.Decode(bits)
		if err2 == nil {
			return code, nil
		} else {
			err2 = fmt.Errorf("keypadDecoder error: %w", err2)
		}
	}
	return code, errors.Join(err1, err2)
}

func (c *CommonDecoder) CheckInputComplete(bits []bool, start time.Time, stop time.Time) bool {
	if c.wiegandDecoder != nil && c.wiegandDecoder.CheckInputComplete(bits, start, stop) {
		return true
	}
	if c.keypadDecoder != nil && c.keypadDecoder.CheckInputComplete(bits, start, stop) {
		return true
	}
	return false
}
