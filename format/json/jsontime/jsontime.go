package jsontime

import (
	"strings"
	"time"
)

// JSONTime takes in account golang time can be an zero/empty and marshals it to an empty string
// instead of "0001-01-01T00:00:00Z" and vice versa.
type JSONTime time.Time

// MarshalJSON marshals a zeop time to an empty string.
func (t JSONTime) MarshalJSON() ([]byte, error) {
	timeVal := time.Time(t)
	if timeVal.IsZero() {
		return []byte(`""`), nil
	}

	return []byte(`"` + timeVal.Format(time.RFC3339Nano) + `"`), nil
}

// UnmarshalJSON unmarshals an empty string and 'mull' to a zero time value.
func (t *JSONTime) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), "`'\" ")

	var timeVal = time.Time{}
	var err error
	if s != "" && s != "null" {
		timeVal, err = time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return err
		}
	}
	*t = JSONTime(timeVal)
	return nil
}

func Now() JSONTime {
	return JSONTime(time.Now().UTC())
}

func (t JSONTime) Time() time.Time {
	return time.Time(t)
}
