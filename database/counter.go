package database

import "time"

// "go.mongodb.org/mongo-driver/v2/bson"
// "go.mongodb.org/mongo-driver/v2/mongo"
// "go.mongodb.org/mongo-driver/v2/mongo/options"

type Reset struct {
	Instigator  string    `bson:"instigator"`
	Reporter    string    `bson:"reporter"`
	Description string    `bson:"description"`
	Timestamp   time.Time `bson:"timestamp"`
}

type Counter struct {
	Id          string  `bson:"_id,omitempty"`
	CreatedBy   string  `bson:"createdBy"`
	Title       string  `bson:"title"`
	Description string  `bson:"description"`
	LastReset   Reset   `bson:"lastReset"`
	History     []Reset `bson:"history"`
}
