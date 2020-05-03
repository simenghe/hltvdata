package main

import (
	"context"
	"fmt"
	"hltvdata/scraper"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// URI mongodb url
var URI = fmt.Sprintf("mongodb+srv://%s:%s@hltvdata-160we.mongodb.net/test?retryWrites=true&w=majority", os.Getenv("MONGODB_USER"), os.Getenv("MONGODB_PASSWORD"))

// UpdateHLTVURLS tries updating the hltv ranking URLs
func UpdateHLTVURLS() time.Duration {
	bench := time.Now()

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

// URLStruct get the urldata
type URLStruct struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	URLS       []string           `bson:"urlList,omitempty"`
	TimeStamp  time.Time          `bson:"timestamp,omitempty"`
	ListLength int                `bson:"listLength,omitempty"`
}

// GetHLTVURLS returns the url list from database
func GetHLTVURLS() (URLStruct, time.Duration) {
	bench := time.Now()
	// URLStruct struct of urllist
	var urlObj URLStruct
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

	opts := options.FindOne()
	opts.SetSort(bson.D{{"timestamp", -1}})
	error := collection.FindOne(ctx, bson.M{}, opts).Decode(&urlObj)
	// cursor, error := collection.Find(ctx, bson.M{"$natural": -1})
	// fmt.Println(cursor)
	fmt.Println(urlObj)
	if error != nil {
		log.Fatal(err)
	}
	// Return the time taken to run this operation.
	return urlObj, time.Since(bench)
}
