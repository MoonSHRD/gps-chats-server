package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullableTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		null, err := json.Marshal(nil)
		return null, err
	}
	unixTimestamp := nt.Time.Unix()
	b, err := json.Marshal(unixTimestamp)
	return b, err
}

func (nt *NullableTime) UnmarshalJSON(b []byte) error {
	var unixTimestamp int64
	err := json.Unmarshal(b, &unixTimestamp)
	*nt = NullableTime{
		Time:  time.Unix(unixTimestamp, 0),
		Valid: true,
	}
	return err
}

// Scan implements the Scanner interface.
func (nt *NullableTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
