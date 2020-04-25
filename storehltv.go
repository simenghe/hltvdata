package main

import (
	"context"
	"fmt"
	"hltvdata/scraper"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// UpdateHLTVURLS tries updating the hltv ranking URLs
func UpdateHLTVURLS() time.Duration {
	bench := time.Now()
	URI := fmt.Sprintf("mongodb+srv://%s:%s@hltvdata-160we.mongodb.net/test?retryWrites=true&w=majority", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PASSWORD"))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		URI,
	))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("hltvdata").Collection("urls")
	urls := scraper.RankingTraverseAsync()
	_, error := collection.InsertOne(ctx, bson.M{
		"urlList":    urls,
		"timestamp":  time.Now(),
		"listLength": len(urls),
	})
	if error != nil {
		log.Fatal(err)
	}
	// Return the time taken to run this operation.
	return time.Since(bench)
}
// GetHLTVURLS returns the url list from database
func GetHLTVURLS(){

}
