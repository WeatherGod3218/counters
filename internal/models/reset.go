package models

import "time"

type Reset struct {
	ResetId     string `json:"reset_id"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}

type ResetListPart struct {
	ResetID          string    `json:"reset_id"`
	ResetDescription string    `json:"reset_description"`
	ResetUsername    string    `json:"reset_username"`
	ResetOccuredAt   time.Time `json:"reset_occured_at"`
}
