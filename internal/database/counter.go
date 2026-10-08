package database

import (
	"context"
	"time"

	"github.com/ComputerScienceHouse/counters/internal/models"
	"github.com/ComputerScienceHouse/counters/internal/util"
	"github.com/jackc/pgx/v5"
)

func CreateCounterWithReset(ctx context.Context, userId string, username string, cReq *models.CreateCounterInput, rReq *models.CreateResetInput) (string, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	counterId, err := GenerateUUID()
	if err != nil {
		return "", err
	}
	rReq.CounterID = counterId

	_, err = tx.Exec(ctx, `
		INSERT INTO counters (counter_id, user_id, username, title, description, last_reset)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, counterId, userId, username, cReq.Title, cReq.Description, nil)
	if err != nil {
		return "", err
	}

	patchedResetTime := util.TranslateTime(rReq.ResetTime)

	resetID, err := CreateResetWithTransaction(ctx, tx, userId, username, patchedResetTime, rReq)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(ctx, `
        UPDATE counters SET last_reset = $1 WHERE counter_id = $2
    `, resetID, counterId)
	if err != nil {
		return "", err
	}

	tx.Commit(ctx)
	return counterId, nil
}

func GetCounters(ctx context.Context) ([]*models.CounterListPart, error) {
	rows, err := db.Query(ctx, `
		SELECT c.counter_id, c.user_id, c.title, c.description, r.description, r.username, r.occured_at FROM counters c
		JOIN resets r ON c.last_reset = r.reset_id
		ORDER BY r.occured_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counters := make([]*models.CounterListPart, 0)

	for rows.Next() {
		var (
			counterId          string
			counterOwnerID     string
			counterTitle       string
			counterDescription string
			resetDescription   string
			resetUsername      string
			resetCreatedAt     time.Time
		)

		if err := rows.Scan(&counterId, &counterOwnerID, &counterTitle, &counterDescription, &resetDescription, &resetUsername, &resetCreatedAt); err != nil {
			return nil, err
		}

		counters = append(counters, &models.CounterListPart{
			CounterID:          counterId,
			CounterOwner:       counterOwnerID,
			CounterTitle:       counterTitle,
			CounterDescription: counterDescription,
			ResetDescription:   resetDescription,
			ResetUsername:      resetUsername,
			ResetOccuredAt:     resetCreatedAt.Unix(),
		})
	}

	return counters, nil
}

func GetCounterFromId(ctx context.Context, rowId string) (*models.CounterListPart, error) {
	var (
		counterId        string
		counterTitle     string
		resetDescription string
		resetUsername    string
		resetOccuredAt   time.Time
	)

	err := db.QueryRow(ctx, `
		SELECT c.counter_id, c.title, r.description, r.username, r.occured_at FROM counters c
		JOIN resets r ON c.last_reset = r.reset_id
		WHERE c.counter_id = $1
	`, rowId).Scan(&counterId, &counterTitle, &resetDescription, &resetUsername, &resetOccuredAt)

	if err != nil {
		return nil, err
	}

	counter := &models.CounterListPart{
		CounterID:        counterId,
		CounterTitle:     counterTitle,
		ResetDescription: resetDescription,
		ResetUsername:    resetUsername,
		ResetOccuredAt:   resetOccuredAt.Unix(),
	}

	return counter, nil
}

func GetCounterOwner(ctx context.Context, counterID string) (string, error) {
	var ownerId string

	_, err := db.Exec(ctx, `
		SELECT user_id FROM counters WHERE counter_id = $1
	`, counterID)

	return ownerId, err
}

func DeleteCounter(ctx context.Context, counterID string) error {
	_, err := db.Exec(ctx, `
		DELETE FROM counters WHERE counter_id = $1
	`, counterID)

	return err
}

func DeleteCounterWithTransaction(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `
		DELETE FROM counters WHERE id = $1
	`, id)

	return err
}

func InitCounters(ctx context.Context) error {
	if _, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS counters (
			counter_id 		UUID PRIMARY KEY,
			user_id 		UUID NOT NULL,
			username 		TEXT NOT NULL,
			title 			TEXT NOT NULL,
			description 	TEXT NOT NULL,
			last_reset	 	UUID,
			created_at 		TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}

	return nil
}

func UpdateCounterLastReset(ctx context.Context, counterId string) (bool, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	var lastReset string

	if err := tx.QueryRow(ctx, `
		SELECT r.reset_id FROM resets r
		WHERE r.counter_id = $1
		ORDER BY r.occured_at DESC
		LIMIT 1
	`, counterId).Scan(&lastReset); err != nil {
		if err == pgx.ErrNoRows {
			return false, DeleteCounter(ctx, counterId)
		}
		return false, err
	}

	if _, err := tx.Exec(ctx, `
		UPDATE counters SET last_reset = $1 WHERE counter_id = $2
	`, lastReset, counterId); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// im doing bullshit. I don't feel like dealing with a migration so this will do
func AlterCounters(ctx context.Context) error {
	if _, err := db.Exec(ctx, `
		DO $$ BEGIN
			ALTER TABLE counters
			ADD CONSTRAINT counters_last_reset_fk
			FOREIGN KEY (last_reset) REFERENCES resets(reset_id)
			ON DELETE SET NULL;		
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$;

	`); err != nil {
		return err
	}

	return nil
}
