package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type ShortDate struct {
	time.Time
}

func (st ShortDate) MarshalJSON() ([]byte, error) {
	if st.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"01-2006"`), nil
}

func (st *ShortDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	t, err := time.Parse(`"01-2006"`, string(data))
	if err != nil {
		return err
	}

	st.Time = t
	return nil
}

func (st ShortDate) Value() (driver.Value, error) {
	if st.IsZero() {
		return nil, nil
	}
	return st.Time, nil
}

func (sd *ShortDate) Scan(value interface{}) error {
	if value == nil {
		sd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		sd.Time = v
		return nil
	case []byte:
		return sd.parse(string(v))
	case string:
		return sd.parse(v)
	default:
		return fmt.Errorf("unknown type: %T", v)
	}
}

func (cd *ShortDate) parse(s string) error {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return fmt.Errorf("wrong date formt '%s': %v", s, err)
	}
	cd.Time = t
	return nil
}
