package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB map[string]any

// Scan implements the sql.Scanner interface. It supports converting from
// string, []byte, or nil into a JSONB value. Attempting to convert from
// any other type will return an error.
func (j *JSONB) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*j = nil
		return nil
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return fmt.Errorf("Scan: unable to scan type %T into JSONB", v)
	}
}

// Value implements the driver.Valuer interface. It converts the JSONB
// value into a SQL driver value which can be used to directly use the
// JSONB as a parameter to a SQL query.
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	// Marshal the JSONB to a byte slice
	return json.Marshal(j)
}

func (j JSONB) String() string {
	// Marshal the JSONB to a byte slice
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}
