package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	var todos []Todo

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(fiber.Map{"todos": todo})
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for _, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todo.Completed = !todo.Completed
				return c.Status(200).JSON(todo)
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
