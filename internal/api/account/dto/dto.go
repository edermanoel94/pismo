package dto

type AccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type AccountResponse struct {
	ID             int    `json:"id,omitempty"`
	DocumentNumber string `json:"document_number"`
}
