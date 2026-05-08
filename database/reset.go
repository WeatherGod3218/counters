package database

import (
	"time"
)

type Reset struct {
	Instigator  string    `bson:"instigator"`
	Reporter    string    `bson:"reporter"`
	Description string    `bson:"description"`
	Timestamp   time.Time `bson:"timestamp"`
}
