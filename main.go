package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"rest/handlers"
	"rest/initializers"
)

func main() {
	app, pool := InitApp()
	defer pool.Close()

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		app.Shutdown()
	}()

	// Start server
	log.Fatal(app.Listen(":8080"))
}

func InitApp() (*fiber.App, *pgxpool.Pool) {
	// Load ENV vars, connect to DB, run migrations
	initializers.Init()

	// Initialize Fiber app
	app := fiber.New()

	// Set task handler
	taskHandler := handlers.NewTaskHandler(initializers.Pool)

	// API setup
	api := app.Group("/api")

	// API Version 1
	v1 := api.Group("/v1")
	{
		// Task handling via task handler
		taskRoute := v1.Group("/tasks")
		{
			// Get ALL tasks
			taskRoute.Get("/", taskHandler.GetAllTasks)

			// CRUD operations on tasks
			taskRoute.Post("/", taskHandler.CreateTask)
			taskRoute.Get("/:id", taskHandler.GetTask)
			taskRoute.Put("/:id", taskHandler.UpdateTask)
			taskRoute.Delete("/:id", taskHandler.DeleteTask)
		}
	}

	return app, initializers.Pool
}
