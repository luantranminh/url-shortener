package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	. "github.com/luantranminh/shorturl/config"
	. "github.com/luantranminh/shorturl/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/speps/go-hashids"
)

var client *mongo.Client
var config Config
var err error

const (
	COLLECTION = "links"
)

// CreateEndpoint .
func CreateEndpoint(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	var request MyURL

	// decode request into MyURL struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	collection := client.Database(config.Database).Collection(COLLECTION)

	//Check if url already exists.
	filter := bson.M{"url": request.URL}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result MyURL
		err := cur.Decode(&result)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"short_url": config.Hostname + "/" + result.ID,
			"url":       result.URL,
		})
		return
	}
	if err := cur.Err(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// create new link
	request.ID = CreateID()

	ctx, c := context.WithTimeout(context.Background(), 60*time.Second)
	_, err = collection.InsertOne(ctx, request)
	defer c()
	if err != nil {
		log.Fatal(err.Error())
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"short_url": config.Hostname + request.ID,
		"url":       request.URL,
	})
}

// Root .
func Root(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	enableCors(&w)
	var request MyURL

	request.ID = params.ByName("id")

	collection := client.Database(config.Database).Collection(COLLECTION)

	//Check if url already exists.
	filter := bson.M{"_id": request.ID}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result MyURL
		err := cur.Decode(&result)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, result.URL, 301)
		return
	}
}

// CreateID creates new ID by hashid with timestamp.
func CreateID() string {
	hd := hashids.NewData()
	hd.Salt = "Guess what?"
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	id, _ := h.Encode([]int{int(now.Unix())})

	return id
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func init() {
	config.Server = os.Getenv("server")
	config.Database = os.Getenv("database")
	config.Hostname = os.Getenv("hostname")
	config.Port = os.Getenv("PORT")

	if config.Server == "" || config.Database == "" {
		config.Read()
	}
}

func main() {

	// connect to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, config.Server)
	if err != nil {
		log.Fatal("Cannot connect to mongoDB by error", err)
	}
	defer cancel()

	router := httprouter.New()
	router.POST("/create", CreateEndpoint)
	router.GET("/:id/", Root)

	if config.Port == "" {
		config.Port = "12345"
	}

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
