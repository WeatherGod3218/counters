package database

import (
	"context"
	"time"

	"github.com/WeatherGod3218/counters/logging"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Reset struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	UserID      string        `bson:"uuid"`
	Reporter    string        `bson:"reporter"`
	Description string        `bson:"description"`
	Timestamp   int64         `bson:"timestamp"`
}

func (r *Reset) IDHex() string {
	return r.Id.Hex()
}

func DeleteReset(ctx context.Context, counterId bson.ObjectID, resetId bson.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := Client.Database(db).Collection("history").UpdateOne(
		ctx,
		bson.M{"_id": counterId},
		bson.M{"$pull": bson.M{
			"history": bson.M{"_id": resetId},
		}},
	)

	if err != nil {
		return false, err
	}

	counter, err := GetCounterFromId(ctx, counterId.Hex())
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "reset", "method": "DeleteReset"}).Warn("error deleting the counter!")
		return false, err
	}

	return counter.Update(ctx)
}
