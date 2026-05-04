package database

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/WeatherGod3218/counters/logging"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type UpsertResult int

const (
	New     UpsertResult = 0
	Updated UpsertResult = 1
)

var Client *mongo.Client
var db = ""

var logger *logrus.Logger = logrus.New()

func Connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("COUNTERS_MONGODB_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "Connect"}).Fatal("error connecting to database, URI: " + uri + "meow")
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "Connect"}).Fatal("error pinging database")
	}

	logging.Logger.WithFields(logrus.Fields{"module": "database", "method": "Connect"}).Info("connected to mongodb")
	db = strings.Split(strings.Split(uri, "/")[3], "?")[0]

	return client
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "Disconnect"}).Fatal("error disconnecting from database")
	}

	logging.Logger.WithFields(logrus.Fields{"module": "database", "method": "Disconnect"}).Info("disconnected from database")
}
