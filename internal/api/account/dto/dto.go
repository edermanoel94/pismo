package dto

type AccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type AccountBalanceRequest struct {
	ID         int     `json:"id"`
	Balance    float64 `json:"balance"`
	NewBalance float64 `json:"new_balance"`
}

type AccountResponse struct {
	ID             int     `json:"id,omitempty"`
	DocumentNumber string  `json:"document_number"`
	Balance        float64 `json:"balance"`
}
