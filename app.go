package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// App contains the router and database connection
type App struct {
	Router *httprouter.Router
	Events *mongo.Collection
}

// Initialize makes the app ready to run
func (a *App) Initialize() {
	a.Router = httprouter.New()
	a.Events = GetCollection("events")
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.GET("/", index)
	a.Router.GET("/events", events)
	a.Router.POST("/event/new", a.createEvent)
}

// Run starts the web server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Hello!")
}

func events(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ALL EVENTS\n")
}

func (a *App) createEvent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var e event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	id, err := e.CreateEvent(a.Events)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"created": id})
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
