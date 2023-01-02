package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"movieApp/model"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://nileshsahitya:<password>@cluster0.yo6exfx.mongodb.net/?retryWrites=true&w=majority"
const dbName = "netflix"
const colName = "watchList"

var collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongodb Connection Success")

	collection = (*mongo.Collection)(client.Database(dbName).Collection((colName)))

	fmt.Println("collection instance is ready")
}

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted 1 movie in db with id: ", inserted.InsertedID)
}

func getAllMovies() []bson.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var movies []bson.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

func updateOneMovie(movieId string) {
	objID, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)

}

func GetAllMyMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)

}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)

	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}
