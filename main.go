// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Task represents a task entity
type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var db *gorm.DB
var err error

func init() {
	// Initialize SQLite database
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the Task model
	db.AutoMigrate(&Task{})
}

func main() {
	app := fiber.New()

	// Routes
	app.Post("/tasks", createTask)
	app.Get("/tasks", getTasks)
	app.Get("/tasks/:id", getTask)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

// Handlers
func createTask(c *fiber.Ctx) error {
	var task Task
	if err := c.BodyParser(&task); err != nil {
		return err
	}

	db.Create(&task)
	return c.JSON(task)
}

func getTasks(c *fiber.Ctx) error {
	var tasks []Task
	db.Find(&tasks)
	return c.JSON(tasks)
}

func getTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.JSON(task)
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	if err := c.BodyParser(&task); err != nil {
		return err
	}

	db.Save(&task)
	return c.JSON(task)
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	db.Delete(&task)
	return c.SendStatus(fiber.StatusNoContent)
}
