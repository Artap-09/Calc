package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	passwordPSQL = os.Getenv("PSQL_PASS")
	namePSQL     = os.Getenv("PSQL_CONTAINER")
	userPSQL     = os.Getenv("PSQL_USER")
	nameMONGO    = os.Getenv("MONGO_CONTAINER")
)

func main() {

	http.HandleFunc("/", handlePastJson)

	err := http.ListenAndServe(":4969", nil)
	if err != nil {
		panic(err)
	}
}

type AB struct {
	A   int `json:"a"`
	B   int `json:"b"`
	Sum int `json:"sum"`
}

func handlePastJson(writer http.ResponseWriter, request *http.Request) {

	var ab AB

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&ab)
	if err != nil {
		panic(err)
	}
	ab.Sum = ab.A + ab.B

	connStr := fmt.Sprintf("user=%s password=%s dbname=mysum sslmode=disable port=5432 host=%s", userPSQL, passwordPSQL, namePSQL)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into calculator values ($1,$2,$3)", ab.A, ab.B, ab.Sum)
	if err != nil {
		panic(err)
	}

	uri := fmt.Sprintf("mongodb://%s:27017", nameMONGO)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = client.Disconnect(context.TODO())

		if err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("mysum").Collection("calculator")

	insertResult, err := collection.InsertOne(context.TODO(), ab)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println(ab)

	io.WriteString(writer, "ok")
}
