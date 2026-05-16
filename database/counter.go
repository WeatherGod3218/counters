package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/WeatherGod3218/counters/logging"
	"github.com/sirupsen/logrus"
)

type Counter struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	UserID      string        `bson:"uuid"`
	CreatedBy   string        `bson:"createdBy"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	LastReset   Reset         `bson:"lastReset"`
}

func CreateCounter(ctx context.Context, counter *Counter) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	newHistory := []Reset{counter.LastReset}

	history := History{
		Id:      counter.Id,
		History: newHistory,
	}

	session, err := Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		_, err := Client.Database(db).Collection("counters").InsertOne(sessCtx, counter)
		if err != nil {
			return nil, err
		}

		_, err = Client.Database(db).Collection("history").InsertOne(sessCtx, history)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func GetCounterFromId(ctx context.Context, id string) (*Counter, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objId, _ := bson.ObjectIDFromHex(id)
	var counter Counter

	if err := Client.Database(db).Collection("counters").FindOne(ctx, bson.M{"_id": objId}).Decode(&counter); err != nil {
		return nil, err
	}

	return &counter, nil
}

func GetAllCounters(ctx context.Context) ([]*Counter, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := Client.Database(db).Collection("counters").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var counters []*Counter
	if err := cursor.All(ctx, &counters); err != nil {
		return nil, err
	}

	return counters, nil
}

func DeleteCounter(ctx context.Context, counterId bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	counter, err := GetCounterFromId(context.Background(), counterId.Hex())
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "reset", "method": "DeleteReset"}).Fatal("error deleting the counter!")
	}

	return counter.Delete(ctx)
}

func (counter *Counter) Delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	session, err := Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		_, err = Client.Database(db).Collection("history").DeleteOne(
			sessCtx,
			bson.M{"_id": counter.Id},
		)

		if err != nil {
			return nil, err
		}

		_, err = Client.Database(db).Collection("counters").DeleteOne(
			sessCtx,
			bson.M{"_id": counter.Id},
		)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "counter", "method": "Delete", "counter": counter.Id}).Warn("error deleting the counter!")
		return err
	}

	return nil
}

func (counter *Counter) Update(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	session, err := Client.StartSession()
	if err != nil {
		return false, err
	}
	defer session.EndSession(context.Background())

	counterEmpty := false

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		var history History

		err = Client.Database(db).Collection("history").FindOne(
			sessCtx,
			bson.M{"_id": counter.Id},
			options.FindOne().SetProjection(bson.M{"history": bson.M{"$slice": 1}}),
		).Decode(&history)

		if err != nil {
			return nil, err
		}

		if len(history.History) == 0 {
			counterEmpty = true
			return nil, nil
		}

		lastReset := history.History[0]
		_, err = Client.Database(db).Collection("counters").UpdateOne(
			sessCtx,
			bson.M{"_id": counter.Id},
			bson.M{"$set": bson.M{"lastReset": lastReset}},
		)
		if err != nil {
			return nil, err
		}

		counter.LastReset = lastReset
		return nil, nil
	})

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "counter", "method": "Update", "counter": counter.Id}).Warn("error updating the counter!")
		return false, err
	}

	if counterEmpty {
		return false, counter.Delete(context.Background())
	}

	return true, nil
}

func (counter *Counter) Reset(ctx context.Context, reset *Reset) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	session, err := Client.StartSession()
	if err != nil {
		return true, err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		_, err = Client.Database(db).Collection("history").UpdateOne(
			sessCtx,
			bson.M{"_id": counter.Id},
			bson.M{"$push": bson.M{
				"history": bson.M{
					"$each": bson.A{reset},
					"$sort": bson.M{"timestamp": -1},
				},
			}},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "counter", "method": "Reset", "counter": counter.Id}).Fatal("error updating counters history!")
	}

	return counter.Update(context.Background())
}

func (c *Counter) IDHex() string {
	return c.Id.Hex()
}
