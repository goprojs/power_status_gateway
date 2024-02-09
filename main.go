package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

type indicator struct {
	ElectricityStatus bool   `json:"estatus"`
	LocationName      string `json:"location_name"`
	LocationID        string `json:"location_id"`
	CurrentTime       string `json:"timestamp"`
}

func main() {
	// Making a Connection to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://arabinda0704:mongos@cluster0.lpazmln.mongodb.net/"))
	if err != nil {
		log.Fatal(err)
	}

	// Set up collection
	collection = client.Database("mgdb").Collection("indicators")

	// Set up Gin router
	router := gin.Default()
	router.GET("/status", getStatus)
	router.POST("/status", postStatus)

	// Start server
	router.Run("localhost:8080")
}

func getStatus(c *gin.Context) {
	var results []indicator
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result indicator
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func postStatus(c *gin.Context) {
	var newInd indicator
	if err := c.BindJSON(&newInd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Inserting data into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newInd)
}
