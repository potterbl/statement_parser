package statement_parser

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/potterbl/statement_parser/pkg/consts"
	"github.com/potterbl/statement_parser/pkg/types"
	"github.com/potterbl/statement_parser/pkg/utils"
	"github.com/sirupsen/logrus"
)

type PrivatBankParser struct{}

func (p *PrivatBankParser) ParsePDF(file multipart.FileHeader, _ []string) ([]types.Transaction, error) {
	csvDataArray, err := utils.ExtractSheetsWithCamelot(&file)
	if err != nil {
		return []types.Transaction{}, err
	}

	var allTransactions []types.Transaction

	for _, csvData := range csvDataArray {
		stringArray, err := utils.ParseCSVToStringArray(csvData)
		if err != nil {
			logrus.Errorf("Failed to parse csv with transactions: %v", err)
			continue
		}

		transactions, err := p.parsePrivatBankCSV(stringArray)
		if err != nil {
			logrus.Errorf("Error parsing PrivatBank CSV: %v", err)
			continue
		}

		allTransactions = append(allTransactions, transactions...)
	}

	return allTransactions, nil
}

func (p *PrivatBankParser) parsePrivatBankCSV(data [][]string) ([]types.Transaction, error) {
	var transactions []types.Transaction

	startRow := 0
	for i, row := range data {
		if len(row) > 0 && strings.Contains(row[0], ".") && len(strings.Split(row[0], ".")) == 3 {
			startRow = i
			break
		}
	}

	for i := startRow; i < len(data); i++ {
		row := data[i]
		if len(row) < 8 {
			continue
		}

		if strings.Contains(row[0], ".") && len(strings.Split(row[0], ".")) == 3 {
			transaction, err := p.parseTransactionRow(data, i)
			if err != nil {
				transaction.ErrorComment = err.Error()
				transactions = append(transactions, transaction)
				continue
			}
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (p *PrivatBankParser) parseTransactionRow(data [][]string, startIndex int) (types.Transaction, error) {
	// PrivatBank CSV format typically spans multiple rows per transaction
	// Row 1: Date, Account, Description, Amount in operation currency, "", "", "", ""
	// Row 2: "", Agreement details, Description, "", Amount in card currency, Commission, Discount, Balance
	// Row 3: Time, Agreement date, Description, Currency, "", "", "", ""

	var transaction types.Transaction

	transaction.Multiplier = 1

	if startIndex >= len(data) || startIndex+2 >= len(data) {
		return types.Transaction{}, fmt.Errorf("insufficient data for transaction")
	}

	row1 := data[startIndex]
	row2 := data[startIndex+1]
	row3 := data[startIndex+2]

	dateParts := []string{strings.TrimSpace(row1[0]), strings.TrimSpace(row3[0])}
	nameParts := []string{strings.TrimSpace(row1[2]), strings.TrimSpace(row2[2]), strings.TrimSpace(row3[2])}

	transaction.Date = strings.Join(dateParts, "T") + ":00Z"
	transaction.Name = strings.TrimSpace(strings.Join(nameParts, " "))

	amount := p.parseAmount(row1[3])
	amountToCompare := p.parseAmount(row2[4])

	currencyStr := strings.TrimSpace(row3[3])
	if currencyAsNum, _ := strconv.Atoi(currencyStr); currencyAsNum != 0 {
		transaction.Currency = consts.CurrencyCodes[currencyAsNum]
	} else {
		transaction.Currency = currencyStr
	}

	if transaction.Currency != "UAH" {
		transaction.Multiplier = amountToCompare / amount
	}

	if amount < 0 {
		transaction.Type = "expense"
	} else {
		transaction.Type = "income"
	}

	transaction.Amount = math.Abs(amount)

	return transaction, nil
}

func (p *PrivatBankParser) parseAmount(amountStr string) float64 {
	cleanAmount := strings.TrimSpace(amountStr)
	cleanAmount = strings.ReplaceAll(cleanAmount, " ", "")
	cleanAmount = strings.ReplaceAll(cleanAmount, ",", ".")

	if cleanAmount == "" {
		return 0.0
	}

	if amount, err := strconv.ParseFloat(cleanAmount, 64); err == nil {
		return amount
	}

	return 0.0
}

func (p *PrivatBankParser) GetBankName() types.BankName {
	return types.BankNamePrivat
}
