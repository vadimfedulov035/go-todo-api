package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Status int

const (
	StatusNew Status = iota // 0 as uninitialized
	StatusInProgress
	StatusDone
)

const StatusConstraint = `Invalid status: must be string from ("new", "in_progress", "done"), got %T %v`

// Stringer interface
func (s Status) String() string {
	switch s {
	case StatusNew:
		return "new"
	case StatusInProgress:
		return "in_progress"
	case StatusDone:
		return "done"
	default:
		return "unknown"
	}
}

// JSON Marshalling (json.Marshaller interface)
func (s Status) MarshalJSON() ([]byte, error) {
	str := s.String()

	// Check value
	if str == "unknown" {
		reason := fmt.Sprintf(StatusConstraint, str, str)
		return nil, NewTaskError(reason)
	}

	return json.Marshal(str)
}

// JSON Unmarshalling (json.Unmarshaller interface)
func (s *Status) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		reason := fmt.Sprintf(StatusConstraint, "non-string", data)
		return NewTaskError(reason)
	}

	// Check value
	switch str {
	case "new":
		*s = StatusNew
	case "in_progress":
		*s = StatusInProgress
	case "done":
		*s = StatusDone
	default:
		reason := fmt.Sprintf(StatusConstraint, str, str)
		return NewTaskError(reason)
	}

	return nil
}

// Database Scanning (sql.Scanner interface)
func (s *Status) Scan(value any) error {
	// Check nilness
	if value == nil {
		reason := fmt.Sprintf(StatusConstraint, "non-string", "NULL")
		return NewTaskError(reason)
	}

	// Check type
	str, ok := value.(string)
	if !ok {
		reason := fmt.Sprintf(StatusConstraint, str, str)
		return NewTaskError(reason)
	}

	// Check value
	switch str {
	case "new":
		*s = StatusNew
	case "in_progress":
		*s = StatusInProgress
	case "done":
		*s = StatusDone
	default:
		reason := fmt.Sprintf(StatusConstraint, str, str)
		return NewTaskError(reason)
	}
	return nil
}

// Database Value (driver.Valuer interface)
func (s Status) Value() (driver.Value, error) {
	str := s.String()

	// Check value
	if str == "unknown" {
		reason := fmt.Sprintf(StatusConstraint, str, str)
		return nil, NewTaskError(reason)
	}

	return str, nil
}
