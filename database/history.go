package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// "go.mongodb.org/mongo-driver/v2/bson"
// "go.mongodb.org/mongo-driver/v2/mongo"
// "go.mongodb.org/mongo-driver/v2/mongo/options"

type History struct {
	Id      bson.ObjectID `bson:"_id,omitempty"`
	History []Reset       `bson:"history"`
}

func CreateHistory(ctx context.Context, history *History) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := Client.Database(db).Collection("history").InsertOne(ctx, history)
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryFromId(ctx context.Context, id string) (*History, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objId, _ := bson.ObjectIDFromHex(id)
	var history History

	if err := Client.Database(db).Collection("history").FindOne(ctx, bson.M{"_id": objId}).Decode(&history); err != nil {
		return nil, err
	}

	return &history, nil
}

// func (counter *Counter) Reset(ctx context.Context, reset *Reset) error {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	objId := counter.Id

// 	_, err := Client.Database(db).Collection("counters").UpdateOne(ctx, bson.M{"_id": objId},
// 		bson.M{
// 			"$set":  bson.M{"lastReset": reset},
// 			"$push": bson.M{"history": reset},
// 		})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (counter *Counter) Close(ctx context.Context) error {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	objId := counter.Id

// 	_, err := Client.Database(db).Collection("counters").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"open": false}})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (c *Counter) IDHex() string {
// 	return c.Id.Hex()
// }

// func GetActiveCounters(ctx context.Context) ([]*Counter, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()

// 	var counters []*Counter

// 	slackThreadReset := Reset{Instigator: "Jeff Mahoney", Reporter: "Eli Mares", Description: "Another thread has hit #ask-eboard", Timestamp: time.Date(2026, 2, 28, 11, 0, 0, 0, time.UTC)}
// 	slackThreadHistory := make([]Reset, 1)
// 	slackThreadHistory[0] = slackThreadReset

// 	counters = append(counters,
// 		&Counter{Id: "1", Title: "Useless Slack Thread", Description: "Hello People, its time", CreatedBy: "weather",
// 			LastReset: slackThreadReset, History: slackThreadHistory},
// 	)

// 	mathExamReset := Reset{Instigator: "Eli Mares", Reporter: "Declan McHale", Description: "Bombed his Calc 2 Final", Timestamp: time.Date(2026, 4, 30, 16, 15, 0, 0, time.UTC)}
// 	mathExamHistory := make([]Reset, 1)
// 	mathExamHistory[0] = mathExamReset

// 	counters = append(counters,
// 		&Counter{Id: "2", Title: "CSHer Failing a Math Class", Description: "Why did you do this", CreatedBy: "weather", LastReset: mathExamReset, History: mathExamHistory},
// 	)

// 	meowReset := Reset{Instigator: "Eli Mares", Reporter: "Noah Opcomm", Description: "Caught his ass meowing in #announcements", Timestamp: time.Date(2026, 5, 3, 12, 51, 0, 0, time.UTC)}
// 	meowResetHistory := make([]Reset, 1)
// 	meowResetHistory[0] = meowReset

// 	counters = append(counters,
// 		&Counter{Id: "3", Title: "meow", Description: "mweoowwww", CreatedBy: "weather", LastReset: meowReset, History: meowResetHistory},
// 	)

// 	igorReset := Reset{Instigator: "Cooper", Reporter: "Eli Mares", Description: "Taking a wild guess all the way from wisconsin", Timestamp: time.Date(2026, 5, 6, 12, 51, 0, 0, time.UTC)}
// 	igorResetHistory := make([]Reset, 1)
// 	igorResetHistory[0] = igorReset

// 	counters = append(counters,
// 		&Counter{Id: "4", Title: "Cooper last played hollow knight", Description: "Hi Cooper!", CreatedBy: "weather", LastReset: igorReset, History: igorResetHistory},
// 	)

// 	return counters, nil
// }
