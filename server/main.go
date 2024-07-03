package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Todo struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongodbUri := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MongodbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection = client.Database("todo_db").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173/",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	Port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + Port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "body is empty"})
	}

	result, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"msg": "insert error"})
	}

	todo.Id = result.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(fiber.Map{"msg": "success"})
}
func updateTodo(c *fiber.Ctx) error {
	Id := c.Params("id")

	ObjectId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "id is invalid"})
	}
	filter := bson.M{"_id": ObjectId}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"msg": "find error"})
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
func deleteTodo(c *fiber.Ctx) error {
	Id := c.Params("id")

	ObjectId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "id is invalid"})
	}
	filter := bson.M{"_id": ObjectId}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"msg": "find error"})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}
