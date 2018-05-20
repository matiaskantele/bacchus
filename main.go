package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type event struct {
	Name   string    `bson:"name"`
	Date   time.Time `bson:"date"`
	Active bool      `bson:"active"`
}

var events *mongo.Collection

func main() {
	events = getCollection("events")
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/events", allEvents)
	router.GET("/event/new", addEvent)

	log.Fatal(http.ListenAndServe(":8080", router))
}
func getCollection(c string) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), os.Getenv("BACCHUS_MONGO_URI"), nil)
	if err != nil {
		log.Fatal("Unable to Connect to database.", err)
	}
	return client.Database("bacchus").Collection(c)
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Hello!")
}

func allEvents(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ALL EVENTS\n")
}

func addEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	urlParams := r.URL.Query()
	name := urlParams.Get("name")
	date, _ := time.Parse(time.RFC3339, urlParams.Get("date"))

	e := event{
		Name:   name,
		Date:   date,
		Active: true,
	}

	_, err := persistToCollection(events, e)
	if err != nil {
		fmt.Fprintf(w, "UNABLE TO ADD EVENT! :(/n")
	}
	fmt.Fprintf(w, "ADDED NEW EVENT!/n")
}

func persistToCollection(c *mongo.Collection, e event) (*mongo.InsertOneResult, error) {
	b, err := bson.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.InsertOne(context.Background(), b)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}
