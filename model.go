package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type event struct {
	ID     objectid.ObjectID `json:"id" bson:"_id"`
	Name   string            `json:"name" bson:"name"`
	Date   time.Time         `json:"date" bson:"date"`
	Active bool              `json:"active" bson:"active"`
}

// Create new event
func (e event) CreateEvent(c *mongo.Collection) (string, error) {
	e.ID = objectid.New()
	e.Active = true
	log.Println("Creating new event:", e)
	return persistToCollection(c, e)
}

// GetCollection returns a collection by name
func GetCollection(c string) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), os.Getenv("BACCHUS_MONGO_URI"), nil)
	if err != nil {
		log.Fatal("Unable to Connect to database.", err)
	}
	return client.Database("bacchus").Collection(c)
}

func persistToCollection(c *mongo.Collection, e event) (string, error) {
	b, err := bson.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.InsertOne(context.Background(), b)
	if err != nil {
		log.Fatal(err)
	}

	var id string
	if oid, ok := res.InsertedID.(*bson.Element); ok {
		id = oid.Value().ObjectID().String()
	}

	return id, err
}
