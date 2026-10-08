package models

type CounterListPart struct {
	CounterID          string `json:"counter_id"`
	CounterOwner       string `json:"counter_owner"`
	CounterTitle       string `json:"counter_title"`
	CounterDescription string `json:"counter_description"`
	ResetDescription   string `json:"reset_description"`
	ResetUsername      string `json:"reset_username"`
	ResetOccuredAt     int64  `json:"reset_occured_at"`
}

type Counter struct {
	CounterID   string `json:"reset_id"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
	LastReset   *Reset `json:"last_reset"`
}
