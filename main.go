package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id, omitempty"`
	Completed bool               `json:"title"`
	Body      string             `json:"body"`
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

	// Create a new collection
	collection = client.Database("todo").Collection("todos") // Create a collection
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Collection created")
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}() // Close the MongoDB connection

	err = client.Ping(ctx, nil) // Ping the MongoDB server
	if err != nil {
		log.Fatal(err) // Handle error
	} else {
		fmt.Println("Pinged MongoDB server")
	}

	app := fiber.New() // Create a new Fiber app

	// ✅ Autoriser toutes les origines (à adapter en prod)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PATCH,DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Get("/api/todos", getTodos)          // Get all todos
	app.Post("/api/todos", createTodo)       // Create a new todo
	app.Patch("/api/todos/:id", updateTodo)  // Update a todo
	app.Delete("/api/todos/:id", deleteTodo) // Delete a todo

	port := os.Getenv("PORT") // Get the port from environment variables
	print("PORT: ", port)     // Print the port
	if port == "" {           // Set a default port if not specified
		port = "5000"
	}
	fmt.Println("Server is running on port: ", port) // Print the port

	// Start the server
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		fmt.Println("Server is ready to accept requests")
	}
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
	}

	if err := cursor.Err(); err != nil { // Check for errors during iteration
		return c.Status(500).SendString("Error iterating through todos") // Handle error
	}

	return c.JSON(todos) // Return the todos as JSON

}

func createTodo(c *fiber.Ctx) error {
	// Create a new todo in the database
	todo := new(Todo)

	// Parse the request body into a todo object
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).SendString("Error parsing todo") // Handle error
	}
	if todo.Body == "" {
		return c.Status(400).SendString("Todo body is required") // Handle error
	}

	// Generate a new ObjectID for the todo
	todo.ID = primitive.NewObjectID()

	// Insert the todo into the database
	_, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(500).SendString("Error inserting todo") // Handle error
	}

	// Return the created todo as JSON
	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	// Update a todo in the database
	// Get the todo ID from the URL parameters
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Convert the ID to an ObjectID
	if err != nil {
		return c.Status(400).SendString("Invalid todo ID") // Handle error
	}
	filter := bson.M{"_id": objectID}                                   // Create a filter to find the todo by ID
	update := bson.M{"$set": bson.M{"completed": true}}                 // Create an update to mark the todo as completed
	_, err = collection.UpdateOne(context.Background(), filter, update) // Update the todo in the database
	if err != nil {
		return c.Status(500).SendString("Error updating todo") // Handle error
	}
	if id == "" {
		return c.Status(400).SendString("Todo ID is required") // Handle error
	}

	return c.Status(200).JSON(fiber.Map{"success": true}) // Return success response
}

func deleteTodo(c *fiber.Ctx) error {
	// Delete a todo from the database
	id := c.Params("id")                           // Get the todo ID from the URL parameters
	objectID, err := primitive.ObjectIDFromHex(id) // Convert the ID to an ObjectID
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"}) // Handle error
	}
	filter := bson.M{"_id": objectID}                           // Create a filter to find the todo by ID
	_, err = collection.DeleteOne(context.Background(), filter) // Delete the todo from the database
	if err != nil {
		return c.Status(500).SendString("Error deleting todo") // Handle error
	}
	return c.Status(200).JSON(fiber.Map{"success delete": true}) // Return success delete response
}
