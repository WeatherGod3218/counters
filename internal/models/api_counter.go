package models

type CreateCounterWithResetInput struct {
	Counter CreateCounterInput `json:"counter"`
	Reset   CreateResetInput   `json:"reset"`
}

type CreateCounterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeleteCounterInput struct {
	RowID string `json:"row_id"`
}
