package database

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Counter struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
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

func (counter *Counter) Reset(ctx context.Context, reset *Reset) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	session, err := Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {

		var history History

		err := Client.Database(db).Collection("history").FindOne(sessCtx, bson.M{"_id": counter.Id}).Decode(&history)
		if err != nil {
			return nil, err
		}

		resetHistory := history.History
		resetHistory = append(resetHistory, *reset)
		sort.Slice(resetHistory, func(i, j int) bool {
			return resetHistory[i].Timestamp.After(resetHistory[j].Timestamp)
		})

		lastReset := resetHistory[0]

		_, err = Client.Database(db).Collection("counters").UpdateOne(
			sessCtx,
			bson.M{"_id": counter.Id},
			bson.M{"$set": bson.M{"lastReset": lastReset}},
		)
		if err != nil {
			return nil, err
		}

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

		counter.LastReset = lastReset

		return nil, nil
	})

	return err
}

func (c *Counter) IDHex() string {
	return c.Id.Hex()
}
