package dto

import (
	"strconv"
	"time"
)

type TimestampTime struct {
	time.Time
}

func (t *TimestampTime) MarshalJSON() ([]byte, error) {
	bin := make([]byte, 16)
	bin = strconv.AppendInt(bin[:0], t.Time.Unix(), 10)
	return bin, nil
}

func (t *TimestampTime) UnmarshalJSON(bin []byte) error {
	v, err := strconv.ParseInt(string(bin), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(v, 0)
	return nil
}
