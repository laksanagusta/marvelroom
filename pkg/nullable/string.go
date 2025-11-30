package nullable

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.
type NullString struct {
	String string
	Valid  bool
}

// NewNullString returns a new NullString from a string pointer
func NewNullString(s *string) NullString {
	if s == nil {
		return NullString{Valid: false}
	}
	return NullString{String: *s, Valid: true}
}

// String returns the string value if valid, otherwise empty string
func (ns NullString) ValueString() string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// StringPtr returns a pointer to the string value if valid, otherwise nil
func (ns NullString) StringPtr() *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	switch v := value.(type) {
	case string:
		ns.String, ns.Valid = v, true
	case []byte:
		ns.String, ns.Valid = string(v), true
	default:
		ns.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// IsSet returns true if the value is set (not null)
func (ns NullString) IsSet() bool {
	return ns.Valid
}

// MarshalJSON implements the json.Marshaler interface
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + ns.String + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		ns.Valid = false
	} else {
		ns.String = *s
		ns.Valid = true
	}
	return nil
}
