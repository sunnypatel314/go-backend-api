package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World!")
	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{
		// {ID: 1, Completed: false, Body: "Learn Python"},
		// {ID: 2, Completed: false, Body: "Learn JavaScript"},
		// {ID: 3, Completed: false, Body: "Learn Go"},
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello world"})
	})

	// read
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(todos)
	})

	// create
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "To-do body is required."})
		}
		todo.ID = 1
		if len(todos) > 0 {
			todo.ID = todos[len(todos)-1].ID + 1
		}
		todos = append(todos, *todo)
		return c.Status(201).JSON(todo)
	})

	// delete
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Could not find todo"})
	})

	// update
	app.Put("/api/todos/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		todo := &Todo{ID: id}
		if err := c.BodyParser(todo); err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
		}
		for i, td := range todos {
			if todo.ID == td.ID {
				todos[i].Body = todo.Body
				todos[i].Completed = todo.Completed
				return c.Status(200).JSON(fiber.Map{"msg": "Successfully updated todo"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
