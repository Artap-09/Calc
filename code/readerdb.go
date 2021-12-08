package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	passwordSQL = os.Getenv("PSQL_PASS")
	nameSQL     = os.Getenv("PSQL_CONTAINER")
	userSQL     = os.Getenv("PSQL_USER")
	nameMONG    = os.Getenv("MONGO_CONTAINER")
)

func main() {
	http.HandleFunc("/", readDB)

	err := http.ListenAndServe(":4979", nil)
	if err != nil {
		panic(err)
	}
}

type DB struct {
	Where string `json:"where"`
	First int    `json:"first"`
	Last  int    `json:"last"`
}

type ABSum struct {
	A   int `json:"a" bson:"a"`
	B   int `json:"b" bson:"b"`
	Sum int `json:"sum" bson:"sum"`
}

func readDB(writer http.ResponseWriter, request *http.Request) {

	var req DB

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}

	switch req.Where {
	case "mongo":
		uri := fmt.Sprintf("mongodb://%s:27017", nameMONG)
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

		rows, err := collection.Find(context.TODO(), bson.D{{"sum", bson.D{{"$gte", req.First}, {"$lt", req.Last}}}})

		var s string

		for rows.Next(context.TODO()) {
			var absum ABSum
			rows.Decode(&absum)
			s = s + fmt.Sprintf("%d + %d = %d\n", absum.A, absum.B, absum.Sum)
		}

		fmt.Print(s)
		io.WriteString(writer, s)
	case "postgres":
		connStr := fmt.Sprintf("user=%s password=%s dbname=mysum sslmode=disable port=5432 host=%s", userSQL, passwordSQL, nameSQL)
		database, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer database.Close()

		rows, err := database.Query("select * from $3 where sum >= $1 and sum < $2", req.First, req.Last, "calculator")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var s string

		for rows.Next() {
			var absum ABSum
			rows.Scan(&absum.A, &absum.B, &absum.Sum)
			s = s + strconv.Itoa(absum.A) + " + " + strconv.Itoa(absum.B) + " = " + strconv.Itoa(absum.Sum) + "\n"
		}

		fmt.Print(s)
		io.WriteString(writer, s)
	default:
		fmt.Println("Неправльно введен тип. Либо \"mongo\", либо \"postgres\"")
		io.WriteString(writer, "Неправльно введен тип. Либо \"mongo\", либо \"postgres\"")
	}
}
