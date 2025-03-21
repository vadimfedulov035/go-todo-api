package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Title string

const TitleConstraint = `Invalid title: must be non-empty string, got %T %q`

// JSON Marshalling (json.Marshaller interface)
func (s Title) MarshalJSON() ([]byte, error) {
	str := string(s)

	// Check value
	if str == "" {
		reason := fmt.Sprintf(TitleConstraint, str, str)
		return nil, NewTaskError(reason)
	}

	return json.Marshal(str)
}

// JSON Unmarshalling (json.Unmarshaller interface)
func (s *Title) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		reason := fmt.Sprintf(TitleConstraint, "non-string", data)
		return NewTaskError(reason)
	}

	// Check value
	if str == "" {
		reason := fmt.Sprintf(TitleConstraint, str, str)
		return NewTaskError(reason)
	}

	*s = Title(str)
	return nil
}

// Database Scanning (sql.Scanner interface)
func (s *Title) Scan(value any) error {
	// Check nilness
	if value == nil {
		reason := fmt.Sprintf(TitleConstraint, "not-string", "NULL")
		return NewTaskError(reason)
	}

	// Check type
	str, ok := value.(string)
	if !ok {
		reason := fmt.Sprintf(TitleConstraint, value, value)
		return NewTaskError(reason)
	}

	// Check value
	if str == "" {
		reason := fmt.Sprintf(TitleConstraint, str, str)
		return NewTaskError(reason)
	}

	*s = Title(str)
	return nil
}

// Database Value (driver.Valuer interface)
func (s Title) Value() (driver.Value, error) {
	str := string(s)

	// Check value
	if str == "" {
		reason := fmt.Sprintf(TitleConstraint, str, str)
		return nil, NewTaskError(reason)
	}

	return str, nil
}
