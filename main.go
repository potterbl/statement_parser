package main

import (
	"bank_parser/pkg/types"
	"bank_parser/pkg/utils/statement_parser"
)

func NewStatementParser(bankName types.BankName) types.BankParser {
	switch bankName {
	case types.BankNameMono:
		return &statement_parser.MonoBankParser{}
	default:
		// Handle unsupported bank names
		return nil
	}
}
