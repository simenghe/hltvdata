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

// GetHLTVURLS returns the url list from database
func GetHLTVURLS() (URLStruct, time.Duration) {
	bench := time.Now()
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

	// Options to setup the query to save the latest one
	opts := options.FindOne()
	opts.SetSort(bson.D{{"timestamp", -1}})
	error := collection.FindOne(ctx, bson.M{}, opts).Decode(&urlObj)
	fmt.Println(urlObj)
	if error != nil {
		log.Fatal(err)
	}
	// Return the time taken to run this operation.
	return urlObj, time.Since(bench)
}

// HLTVRanking stores multiple lists of rankings.
type HLTVRanking struct {
	CSGOTeams []scraper.CSGOteam
	Timestamp time.Time
	URL       string
}

// UpdateHLTVRankings lists off a range of rankings through time
func UpdateHLTVRankings() bson.M {
	// Grab HLTVURLS
	urlObj, _ := GetHLTVURLS()
	c := make(chan []scraper.CSGOteam)
	urlCount := 0
	for _, s := range urlObj.URLS {
		go scraper.ScrapeHltvTeamsByURLAsync(s, c)
		// time.Sleep(time.Millisecond * 80)
		urlCount++
	}
	HLTVRankingCollection := make([]HLTVRanking, urlCount)
	for i, s := range urlObj.URLS {
		HLTVRankingCollection[i] = HLTVRanking{
			CSGOTeams: <-c, // exhaust the channel
			Timestamp: time.Now(),
			URL:       s,
		}
	}
	// Insert Into The Database
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
	collection := client.Database("hltvdata").Collection("rankings")
	HLTVcollection := bson.M{
		"timestamp":  time.Now(),
		"urlCount":   urlCount,
		"collection": HLTVRankingCollection,
	}
	_, error := collection.InsertOne(ctx, HLTVcollection)
	if error != nil {
		log.Fatal(err)
	}
	return HLTVcollection
}

// GetHLTVRankings grabs the latest collection from mongo.
func GetHLTVRankings() bson.M {
	var rankCollection RankingStruct
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
	collection := client.Database("hltvdata").Collection("rankings")

	// Options to setup the query to save the latest one
	opts := options.FindOne()
	opts.SetSort(bson.D{{"timestamp", -1}})
	error := collection.FindOne(ctx, bson.M{}, opts).Decode(&rankCollection)
	returnObject := bson.M{
		"ID":         rankCollection.ID,
		"timestamp":  rankCollection.TimeStamp,
		"urlcount":   rankCollection.URLCount,
		"collection": rankCollection.Collection,
	}
	if error != nil {
		log.Fatal(err)
	}
	// Return the time taken to run this operation.
	return returnObject
}
