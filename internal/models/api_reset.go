package models

type CreateResetInput struct {
	CounterID   string `json:"counter_id"`
	Description string `json:"description"`
	ResetTime   int64  `json:"reset_time"`
}

type DeleteResetInput struct {
	RowID string `json:"reset_id"`
}
