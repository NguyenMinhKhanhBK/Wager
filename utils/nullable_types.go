package utils

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint struct {
	Uint  uint
	Valid bool // Valid is true if Uint is not NULL
}

func NewNullUint(value uint) NullUint {
	return NullUint{Uint: value, Valid: true}
}

// Scan implements the Scanner interface.
func (n *NullUint) Scan(value interface{}) error {
	if value == nil {
		n.Uint, n.Valid = 0, false
		return nil
	}

	var temp sql.NullInt64
	if err := temp.Scan(value); err != nil {
		return err
	}
	n.Uint = uint(temp.Int64)
	n.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint), nil
}

func (n NullUint) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Uint)
}

func (n *NullUint) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Uint)
	n.Valid = (err == nil)
	return err
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(nf.Float64)
}

func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}
