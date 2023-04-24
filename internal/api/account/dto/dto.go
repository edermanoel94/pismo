package dto

type AccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type AccountResponse struct {
	ID             int    `json:"id,omitempty"`
	DocumentNumber string `json:"document_number"`
}
