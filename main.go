package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sensor struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	sensorID   int                `bson:"sensorId,omitempty"`
	numSamples int                `bson:"numSamples,omitempty"`
	day        string             `bson:"day,omitempty"`
	first      float64            `bson:"first,omitempty"`
	last       float64            `bson:"last,omitempty"`
	samples    []sample           `bson:"samples,omitempty"`
}

type sample struct {
	val  float64 `bson:"val,omitempty"`
	time int64   `bson:"time,omitempty"`
}

func main() {
	/*ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()*/
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://EASC_01:0KzQlxBvEk8yDmLV@thundersandbox.v0ydj.mongodb.net/ThunderSandBox?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	bmp280DB := client.Database("bmp280")
	locations := bmp280DB.Collection("locations")
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Sensor")
	}).Methods("GET")
	r.HandleFunc("/sensor/{sensorID}/temp/{temp}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sensorID, _ := strconv.Atoi(vars["sensorID"])
		temp, _ := strconv.ParseFloat(vars["temp"], 64)

		fmt.Fprintf(w, "pushed sensor: %d with temp %f\n", sensorID, temp)
		opts := options.Update().SetUpsert(true)
		filter := bson.D{{Key: "sensorID", Value: sensorID}, {Key: "numSamples", Value: bson.D{{Key: "$lt", Value: 250}}}}
		update := bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "samples", Value: bson.D{
					{Key: "val", Value: temp},
					{Key: "time", Value: time.Now().UnixNano()},
				}}},
			},
			{Key: "$inc", Value: bson.D{{Key: "numSamples", Value: 1}}},
		}
		locationUpdate, err := locations.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(locationUpdate)
	}).Methods("GET")

	r.HandleFunc("/sensor/{sensorID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sensorID, _ := strconv.Atoi(vars["sensorID"])
		var values []bson.M
		cursor, err := locations.Find(ctx, bson.M{"sensorID": bson.D{{"$eq", sensorID}}})
		if err != nil {
			panic(err)
		}
		if err = cursor.All(ctx, &values); err != nil {
			panic(err)
		}
		for _, value := range values {
			fmt.Println(value["numSamples"])
		}

		fmt.Fprintf(w, "Reading MongoDB\n")
	}).Methods("GET")

	r.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		sensor := sensor{
			sensorID:   1,
			numSamples: 250,
			day:        time.Now().Format("01-02-2006"),
			first:      0,
			last:       0,
			samples:    []sample{{0, 0}},
		}
		insertResult, err := locations.InsertOne(ctx, sensor)
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult.InsertedID)
	}).Methods("GET")

	http.ListenAndServe(":"+port, r)
}
