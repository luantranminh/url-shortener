package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

const (
	COLLECTION = "links"
)

// CreateEndpoint .
func CreateEndpoint(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request MyURL

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

}

// ExpandEndpoint .
func ExpandEndpoint(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Root .
func Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello and welcome to the url shortener service")
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

func main() {
	config.Read()

	// connect to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, config.Server)
	if err != nil {
		log.Fatal("Cannot connect to mongoDB by error", err)
	}
	defer cancel()

	collection := client.Database(config.Database).Collection(COLLECTION)
	ctx, c := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159, "_id": "asdsax"})
	defer c()
	fmt.Println(bson.RawValue{})
	if err != nil {
		log.Fatal(err.Error())
	}

	id := res.InsertedID
	fmt.Print(id)

	router := httprouter.New()
	router.POST("/create", CreateEndpoint)
	router.GET("/:id", ExpandEndpoint)
	router.GET("/", Root)

	log.Fatal(http.ListenAndServe(":12345", router))
}
