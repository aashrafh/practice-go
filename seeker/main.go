package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	Id          primitive.ObjectID `json:"_id,omitempty bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty bson:"title,omitempty"`
	Description string             `json:"description,omitempty bson:"description,omitempty"`
	Image       string             `json:"image,omitempty bson:"image,omitempty"`
	Price       int                `json:"price,omitempty bson:"price,omitempty"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("seeker").Collection("products")

	app := fiber.New()

	app.Use(cors.New())

	app.Listen(":3000")

}
