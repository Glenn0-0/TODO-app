package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type todo struct {
	Item string
}

func main() {
	connStr := "postgresql://postgres:111111@localhost/todos?sslmode=disable"

	// connect to the database.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5432"
	}

	app.Static("/", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var todos []string

	rows, err := db.Query("SELECT * FROM TODOs")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occurred")
	}

	for rows.Next() {
		rows.Scan(&res)
		todos = append(todos, res)
	}

	return c.Render("index", fiber.Map {
		"TODOs": todos,
	})
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error occurred: %v", err)
		return c.SendString(err.Error())
	}

	fmt.Printf("%v", newTodo)

	if newTodo.Item != "" {
		_, err := db.Exec("INSERT into TODOs VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatal("An error occurred while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	olditem := c.Query("olditem")
	newitem := c.Query("newitem")

	db.Exec("UPDATE TODOs SET item=$1 WHERE item=$2", newitem, olditem)

	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	todoToDelete := c.Query("item")

	db.Exec("DELETE from TODOs WHERE item=$1", todoToDelete)
	
	return c.SendString("deleted")
}