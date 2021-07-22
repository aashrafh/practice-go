package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	app := fiber.New()

	app.Use(cors.New())

	app.Post("/api/products/populate", func(c *fiber.Ctx) error {
		collection := client.Database("seeker").Collection("products")
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for i := 0; i < 50; i++ {
			collection.InsertOne(ctx, Product{
				Title:       faker.Word(),
				Description: faker.Paragraph(),
				Image:       fmt.Sprintf("http://lorempixel.com/200/200?%s", faker.UUIDDigit()),
				Price:       rand.Intn(90) + 10,
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	app.Get("/api/products/frontend", func(c *fiber.Ctx) error {
		collection := client.Database("seeker").Collection("products")
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var products []Product

		cur, err := collection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var product Product

			err := cur.Decode(&product)
			if err != nil {
				log.Fatal(err)
			}

			products = append(products, product)
		}

		return c.JSON(products)
	})

	app.Get("/api/products/backend", func(c *fiber.Ctx) error {
		collection := client.Database("seeker").Collection("products")
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var products []Product

		filter := bson.M{}
		findOptions := options.Find()

		if search := c.Query("search"); search != "" {
			filter = bson.M{
				"$or": []bson.M{
					{
						"title": bson.M{
							"$regex": primitive.Regex{
								Pattern: search,
								Options: "i",
							},
						},
					},
					{
						"description": bson.M{
							"$regex": primitive.Regex{
								Pattern: search,
								Options: "i",
							},
						},
					},
				},
			}
		}

		if sort := c.Query("sort"); sort != "" {
			if sort == "asc" {
				findOptions.SetSort(bson.D{{"price", 1}})
			} else if sort == "desc" {
				findOptions.SetSort(bson.D{{"price", -1}})
			}
		}

		page, _ := strconv.Atoi(c.Query("page", "1"))
		var perPage int64 = 9
		total, _ := collection.CountDocuments(ctx, filter)
		findOptions.SetSkip((int64(page) - 1) * perPage)
		findOptions.SetLimit(perPage)

		cur, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var product Product

			err := cur.Decode(&product)
			if err != nil {
				log.Fatal(err)
			}

			products = append(products, product)
		}

		return c.JSON(fiber.Map{
			"data":      products,
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(total) / float64(perPage)),
		})
	})

	app.Listen(":3000")

}
