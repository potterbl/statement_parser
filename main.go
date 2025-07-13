package statement_parser

import (
	"github.com/potterbl/statement_parser/pkg/types"
	"github.com/potterbl/statement_parser/pkg/utils/statement_parser"
)

func NewStatementParser(bankName types.BankName) types.BankParser {
	switch bankName {
	case types.BankNameMono:
		return &statement_parser.MonoBankParser{}
	case types.BankNamePrivat:
		return &statement_parser.PrivatBankParser{}
	default:
		// Handle unsupported bank names
		return nil
	}
}
