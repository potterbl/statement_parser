package statement_parser

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"

	"github.com/potterbl/statement_parser/pkg/types"
	"github.com/potterbl/statement_parser/pkg/utils"
)

type MonoBankParser struct{}

func (p *MonoBankParser) ParsePDF(file multipart.FileHeader, languages []string) ([]types.Transaction, error) {
	fileChunks, err := utils.ParsePDFFileWithOCR(&file, languages)
	if err != nil {
		return []types.Transaction{}, err
	}

	var (
		wg           sync.WaitGroup
		transactions []types.Transaction
		allErrors    []error
		mu           sync.Mutex
	)

	for _, chunk := range fileChunks {
		wg.Add(1)
		go func() {
			mu.Lock()

			defer wg.Done()
			defer mu.Unlock()

			trx, err := p.parseStatement(chunk)

			if err != nil {
				allErrors = append(allErrors, err)
				return
			}

			transactions = append(transactions, trx...)
		}()
	}

	wg.Wait()

	if len(allErrors) > 0 {
		return transactions, fmt.Errorf("errors occurred during parsing: %v", allErrors)
	}

	return transactions, nil
}

func (p *MonoBankParser) parseStatement(ocrText string) ([]types.Transaction, error) {
	var transactions []types.Transaction
	lines := strings.Split(ocrText, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || (!strings.Contains(line, ".2025") && !strings.Contains(line, ".2024") && !strings.Contains(line, ".2023")) {
			continue
		}

		// Пример строки:
		// 26.06.2025 Migros 5411 -408.27 -385.85 TRY 1.06 0.00 4.08 1041.94
		tokens := strings.Fields(line)
		if len(tokens) < 7 {
			continue
		}

		dateStr := tokens[0]
		name := tokens[1]
		amountStr := tokens[3]
		currency := tokens[5]

		dateTime := fmt.Sprintf("%sT00:00:00Z", dateStr)
		if i+1 < len(lines) && strings.Contains(lines[i+1], ":") {
			timeLine := strings.TrimSpace(lines[i+1])
			dateTime = dateStr + "T" + strings.ReplaceAll(timeLine, " ", "") + "Z"
			i++
		}

		amount, _ := strconv.ParseFloat(strings.ReplaceAll(amountStr, ",", "."), 64)

		t := types.Transaction{
			Amount:   math.Abs(amount),
			Currency: currency,
			Name:     name,
			Type:     "expense",
			Date:     dateTime,
		}

		if amount > 0 {
			t.Type = "income"
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (m *MonoBankParser) GetBankName() types.BankName {
	return types.BankNameMono
}
