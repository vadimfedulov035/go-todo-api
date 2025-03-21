package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rest/models"
)

type TaskHandler struct {
	Pool *pgxpool.Pool
}

func NewTaskHandler(pool *pgxpool.Pool) *TaskHandler {
	return &TaskHandler{Pool: pool}
}

func (t *TaskHandler) GetAllTasks(c fiber.Ctx) error {
	// SQL query: get all tasks
	const query = `SELECT
		id, title, description, status, created_at, updated_at
		FROM private.tasks
	`

	// Get all rows
	rows, _ := t.Pool.Query(c.Context(), query)
	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to process tasks",
			"reason": models.GetErrorReason(err),
		})
	}

	log.Print("Retrieved all tasks")
	return c.JSON(tasks)
}

func (t *TaskHandler) CreateTask(c fiber.Ctx) error {
	// Parse request body
	task := new(models.Task)
	err := c.Bind().Body(&task)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid request body",
			"reason": models.GetErrorReason(err),
		})
	}

	// SQL query: create task
	const query = `
        INSERT INTO private.tasks (title, description, status)
        VALUES ($1, $2, $3)
        RETURNING id, title, description, status, created_at, updated_at
    `

	// Execute query
	err = t.Pool.QueryRow(c.Context(), query,
		task.Title,
		task.Description,
		task.Status,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to create task",
			"reason": models.GetErrorReason(err),
		})
	}

	log.Printf("Task %d created", task.ID)
	return c.Status(fiber.StatusCreated).JSON(task)
}

func (t *TaskHandler) GetTask(c fiber.Ctx) error {
	// Get passed ID
	id := c.Params("id")

	// SQL query: get task by ID
	const query = `
		SELECT id, title, description, status, created_at, updated_at
		FROM private.tasks
		WHERE id = $1
	`

	// Execute query
	task := new(models.Task)
	err := t.Pool.QueryRow(c.Context(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Failed to get task",
			"reason": models.GetErrorReason(err),
		})
	}

	log.Printf("Task %d retrieved", task.ID)
	return c.JSON(task)
}

func (t *TaskHandler) UpdateTask(c fiber.Ctx) error {
	// Get passed ID
	id, err := strconv.Atoi(c.Params("id"))

	// Log and expose ID error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Failed to update task",
			"reason": "Invalid task ID format",
		})
	}

	// Parse request body
	var task models.Task
	err = c.Bind().Body(&task)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid request body",
			"reason": models.GetErrorReason(err),
		})
	}

	// SQL query: update task by ID (created_at stays)
	const query = `
		UPDATE private.tasks
		SET title = $2, description = $3, status = $4, updated_at = $5
		WHERE ID = $1
		RETURNING id, title, description, status, created_at, updated_at
	`

	// Execute query
	task.ID = id
	task.UpdatedAt = time.Now()
	err = t.Pool.QueryRow(c.Context(), query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to update task",
			"reason": models.GetErrorReason(err),
		})
	}

	log.Printf("Task %d updated", task.ID)
	return c.JSON(task)
}

func (t *TaskHandler) DeleteTask(c fiber.Ctx) error {
	// Get passed ID
	id, err := strconv.Atoi(c.Params("id"))

	// Log and expose ID error
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Failed to delete task",
			"reason": "Invalid task ID format",
		})
	}

	// SQL query: delete task by ID
	const query = `
		DELETE FROM private.tasks
		WHERE id = $1
		RETURNING id, title, description, status, created_at, updated_at
	`

	// Execute query
	task := new(models.Task)
	err = t.Pool.QueryRow(c.Context(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	// Log any error, expose only safe errors
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to delete task",
			"reason": models.GetErrorReason(err),
		})
	}

	log.Printf("Task %d deleted", task.ID)
	return c.JSON(task)
}
