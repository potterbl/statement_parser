package types

import "mime/multipart"

type BankParser interface {
	ParsePDF(file multipart.FileHeader, languages []string) ([]Transaction, error)
	GetBankName() BankName
}

type Transaction struct {
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Name     string  `json:"name"`
	Date     string  `json:"date"`
}

type BankName string

const (
	BankNameMono BankName = "Mono"
)
