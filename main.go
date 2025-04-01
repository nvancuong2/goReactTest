package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        int    `json:"id" bson:"_id"`
	Completed bool   `json:"title"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env") // Load environment variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}

	Mongo_URI := os.Getenv("MONGO_URI") // Get the MongoDB URI from environment variables

	if Mongo_URI == "" {
		fmt.Println("MONGO_URI not set")
	} else {
		fmt.Println("MONGO_URI: ", Mongo_URI)
	}

	// Import the MongoDB driver
	var ctx = context.Background()                           // MongoDB connection string
	var clientOptions = options.Client().ApplyURI(Mongo_URI) // Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)         // Check the connection

	if err != nil {
		log.Fatal(err) // Handle error
	} else {
		fmt.Println("Connected to MongoDB")
	}
	err = client.Ping(ctx, nil) // Ping the MongoDB server
	if err != nil {
		log.Fatal(err) // Handle error
	} else {
		fmt.Println("Pinged MongoDB server")
	}
	// Create a new collection

	// collection = client.Database("test").Collection("todos") // Create a new collection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}() // Close the MongoDB connection

	// Check the connection
	collection = client.Database("todo").Collection("todos") // Create a collection
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Collection created")
	}

	app := fiber.New()                       // Create a new Fiber app
	app.Get("/api/todos", getTodos)          // Get all todos
	app.Post("/api/todos", createTodo)       // Create a new todo
	app.Patch("/api/todos/:id", updateTodo)  // Update a todo
	app.Delete("/api/todos/:id", deleteTodo) // Delete a todo

	port := os.Getenv("PORT") // Get the port from environment variables
	if port == "" {           // Set a default port if not specified
		port = "5000"
	}
	fmt.Println("Server is running on port: ", port) // Print the port
}
func getTodos(c *fiber.Ctx) error {
	// Get all todos from the database
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{}) // Find all todos in the collection

	if err != nil {
		return c.Status(500).SendString("Error fetching todos") // Handle error
	}

	for cursor.Next(context.Background()) { // Iterate through the cursor
		var todo Todo
		if err := cursor.Decode(&todo); err != nil { // Decode the cursor into a todo object
			return c.Status(500).SendString("Error decoding todo") // Handle error
		}
		todos = append(todos, todo) // Append the todo to the slice
		return c.JSON(todos)        // Return the todos as JSON
	}
	return nil // Ensure the function ends properly
}

func createTodo(c *fiber.Ctx) error {
	// Create a new todo in the database
	return c.SendString("Create a new todo")
}
func updateTodo(c *fiber.Ctx) error {
	// Update a todo in the database
	return c.SendString("Update a todo")
}
func deleteTodo(c *fiber.Ctx) error {
	// Delete a todo from the database
	return c.SendString("Delete a todo")
}
