package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

// Main model
type Task struct {
	ID          int       `json:"id"`
	Title       Title     `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Task error exposed to end user
type TaskError struct {
	Reason string `json:"reason"`
}

// Error implementation
func (c *TaskError) Error() string {
	return fmt.Sprintf("%s", c.Reason)
}

// Constructor
func NewTaskError(reason string) *TaskError {
	return &TaskError{
		Reason: reason,
	}
}

// Expose error reason to end user
func GetErrorReason(err error) string {
	// Task error
	var taskErr *TaskError
	if errors.As(err, &taskErr) {
		return taskErr.Error()
	}

	// General JSON error
	if errors.Is(err, fiber.ErrUnprocessableEntity) {
		return err.Error()
	}

	// General NotFound error
	if errors.Is(err, pgx.ErrNoRows) {
		return err.Error()
	}

	return "Something went wrong"
}
