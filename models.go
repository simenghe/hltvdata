package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URLStruct get the urldata
type URLStruct struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	URLS       []string           `bson:"urlList,omitempty"`
	TimeStamp  time.Time          `bson:"timestamp,omitempty"`
	ListLength int                `bson:"listLength,omitempty"`
}

// CSGOteam as defined as before.
type CSGOteam struct {
	TeamName   string   `bson:"teamname,omitempty"`
	Points     int      `bson:"points,omitempty"` // Points need to be int
	Ranking    int      `bson:"ranking,omitempty"`
	Date       string   `bson:"date,omitempty"`
	PlayerList []string `bson:"playerlist,omitempty"`
}

// RankingStruct for the collections of CSGO rankings
type RankingStruct struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	URLCount   int                `bson:"urlCount,omitempty"`
	TimeStamp  time.Time          `bson:"timestamp,omitempty"`
	Collection []bson.M           `bson:"collection,omitempty"` // The good stuff
}
