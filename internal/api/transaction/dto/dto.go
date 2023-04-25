package dto

type TransactionRequest struct {
	AccountId       int     `json:"account_id" validate:"required"`
	OperationTypeId int     `json:"operation_type_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
}

type TransactionResponse struct {
	Amount float64 `json:"amount"`
}
