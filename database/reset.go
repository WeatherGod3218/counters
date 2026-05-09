package database

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Reset struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	Instigator  string        `bson:"instigator"`
	Reporter    string        `bson:"reporter"`
	Description string        `bson:"description"`
	Timestamp   time.Time     `bson:"timestamp"`
}
