package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kpunith8/go-jwt-auth/utils"
)

// DBConnection - DB Connection singleton
var DBConnection = Connect()

// Connect - Function to connect to mongoDB
func Connect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(utils.GetEnvVariable("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
	// Ping to verify DB is connected
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Available databases", databases)

	// testDatabase := client.Database("test")
	// goUserCollection := testDatabase.Collection("go-user")

	// Don't insert the same data multiple times, need to fix this extracting this to a function
	// goUserResult, err := goUserCollection.InsertOne(ctx, bson.D{
	// 	{"name", "Punith K"},
	// 	{"age", 20},
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted documents into go-user collection!\n", goUserResult)

	// Using struct models
	// usersCursor, err := goUserCollection.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// modelGoUsers := []models.User{}

	// if err = usersCursor.All(ctx, &modelGoUsers); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("All users", modelGoUsers)

	// Using bson model
	// var goUser bson.M
	// if err = goUserCollection.FindOne(ctx, bson.M{}).Decode(&goUser); err != nil {
	// 	log.Fatal(err)
	// }

	// using struct
	// modelGoSingleUser := models.User{}
	// if err = goUserCollection.FindOne(ctx, bson.M{}).Decode(&modelGoSingleUser); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("A single user\n", modelGoSingleUser)
}
