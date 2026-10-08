package database

import (
	"context"
	"time"

	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/jackc/pgx/v5"
)

func CreateReset(ctx context.Context, userID string, username string, resetTime int64, req *models.CreateResetInput) (string, error) {
	resetID, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO resets (reset_id, counter_id, user_id, username, description, occured_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, resetID, req.CounterID, userID, username, req.Description, time.Unix(resetTime, 0))

	return resetID, err
}

func CreateResetWithTransaction(ctx context.Context, tx pgx.Tx, userID string, username string, resetTime int64, req *models.CreateResetInput) (string, error) {
	resetID, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO resets (reset_id, counter_id, user_id, username, description, occured_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, resetID, req.CounterID, userID, username, req.Description, time.Unix(resetTime, 0))

	return resetID, err
}

func DeleteReset(ctx context.Context, counterid string, resetid string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = db.Exec(ctx, `
		DELETE FROM resets WHERE reset_id = $1
	`, resetid)

	return err
}

func GetResetOwner(ctx context.Context, resetID string) (string, error) {
	var ownerId string

	err := db.QueryRow(ctx, `
		SELECT user_id FROM resets WHERE reset_id = $1
	`, resetID).Scan(&ownerId)

	return ownerId, err
}

func GetCounterFromReset(ctx context.Context, resetID string) (string, error) {
	var counterID string

	if err := db.QueryRow(ctx, `
		SELECT r.counter_id FROM resets r WHERE r.reset_id = $1
	`, resetID).Scan(&counterID); err != nil {
		return "", err
	}

	return counterID, nil
}

func GetResetsFromCounterId(ctx context.Context, counterID string) ([]*models.ResetListPart, error) {
	rows, err := db.Query(ctx, `
		SELECT r.reset_id, r.description, r.user_id, r.username, r.occured_at FROM resets r
		WHERE r.counter_id = $1
		ORDER BY r.occured_at DESC
	`, counterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resets := make([]*models.ResetListPart, 0)

	for rows.Next() {
		var (
			resetId          string
			resetDescription string
			resetOwner       string
			resetUsername    string
			resetOccuredAt   time.Time
		)

		if err := rows.Scan(&resetId, &resetDescription, &resetOwner, &resetUsername, &resetOccuredAt); err != nil {
			return nil, err
		}

		resets = append(resets, &models.ResetListPart{
			ResetID:          resetId,
			ResetDescription: resetDescription,
			ResetOwner:       resetOwner,
			ResetUsername:    resetUsername,
			ResetOccuredAt:   resetOccuredAt.Unix(),
		})
	}

	return resets, nil
}

func InitReset(ctx context.Context) error {
	if _, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS resets (
			reset_id 	UUID PRIMARY KEY,
			counter_id 	UUID REFERENCES counters(counter_id) ON DELETE CASCADE,
			user_id 	UUID NOT NULL,
			username 	TEXT NOT NULL,
			description TEXT NOT NULL,
			occured_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}

	return nil
}
