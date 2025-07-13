package statement_parser

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"

	"github.com/potterbl/statement_parser/pkg/consts"
	"github.com/potterbl/statement_parser/pkg/types"
	"github.com/potterbl/statement_parser/pkg/utils"
)

const (
	epsilon = 30
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

		line = strings.ReplaceAll(line, "|", "")

		tokens := strings.Fields(line)

		if len(tokens) < 8 {
			continue
		}

		dateStr := tokens[0]

		var currencyIndex = -1
		for j := 1; j < len(tokens); j++ {
			for _, c := range consts.Currencies {
				if tokens[j] == c {
					currencyIndex = j
					break
				}
			}
			if currencyIndex != -1 {
				break
			}
		}

		if currencyIndex == -1 {
			continue
		}

		amount, multiplier, currency, errorComment, mccIndex := tryParseAmount(tokens, currencyIndex)

		var nameTokens []string
		if mccIndex > 1 {
			nameTokens = tokens[1:mccIndex]
		} else {
			nameTokens = tokens[1:2]
		}

		dateTime := fmt.Sprintf("%sT00:00:00Z", dateStr)
		if i+1 < len(lines) && strings.Contains(lines[i+1], ":") {
			timeLine := strings.TrimSpace(lines[i+1])
			timeLine = strings.ReplaceAll(timeLine, "|", "")

			time := strings.Fields(timeLine)

			nameTokens = append(nameTokens, strings.Join(time[1:], " "))

			dateTime = dateStr + "T" + time[0] + "Z"
			i++
		}

		name := strings.Join(nameTokens, " ")

		t := types.Transaction{
			Amount:       math.Abs(amount),
			Currency:     currency,
			Multiplier:   multiplier,
			Name:         name,
			Type:         "expense",
			Date:         dateTime,
			ErrorComment: errorComment,
		}

		if amount > 0 {
			t.Type = "income"
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func tryParseAmount(tokens []string, currencyIndex int) (float64, float64, string, string, int) {
	currency := tokens[currencyIndex]

	multiplier := 1.0
	if currency != "UAH" && currencyIndex+1 < len(tokens) {
		multiplierStr := strings.ReplaceAll(tokens[currencyIndex+1], "—", "")
		multiplierStr = strings.ReplaceAll(multiplierStr, ",", ".")
		multiplier, _ = strconv.ParseFloat(multiplierStr, 64)
	}

	mccCodeIndex := -1
	for i := currencyIndex; i >= 0; i-- {
		parsedNum, _ := strconv.ParseFloat(tokens[i], 64)
		if len(tokens[i]) == 4 && !math.IsNaN(parsedNum) && parsedNum != 0.0 {
			mccCodeIndex = i
		}
	}

	if mccCodeIndex == -1 {
		return 0.0, 0.0, currency, fmt.Sprintf("No MCC code found, raw: %s", strings.Join(tokens, " ")), -1
	}

	amountToCompare := 0.0
	amount := 0.0

	secondAmountFirstIndex := -1
	for i := mccCodeIndex + 1; i < currencyIndex; i++ {
		if strings.Contains(tokens[i], ".") {
			if secondAmountFirstIndex == -1 {
				secondAmountFirstIndex = i + 1
				break
			}
		}
	}
	if secondAmountFirstIndex == -1 {
		return 0.0, 0.0, currency, fmt.Sprintf("Not found end of first amount, raw: %s", strings.Join(tokens, " ")), -1
	}

	amountToCompareStr := cleanAmountString(strings.Join(tokens[mccCodeIndex+1:secondAmountFirstIndex], ""))
	amountStr := cleanAmountString(strings.Join(tokens[secondAmountFirstIndex:currencyIndex], ""))

	amountToCompare, _ = strconv.ParseFloat(strings.ReplaceAll(amountToCompareStr, ",", "."), 64)
	amount, _ = strconv.ParseFloat(strings.ReplaceAll(amountStr, ",", "."), 64)

	if math.Abs(amount*multiplier)-math.Abs(amountToCompare) < epsilon {
		return amount, multiplier, currency, "", mccCodeIndex
	}

	return amount, multiplier, currency, fmt.Sprintf("Amount with multiplier does not match the default amount in currency, raw: %s", strings.Join(tokens, " ")), mccCodeIndex
}

func cleanAmountString(s string) string {
	var result strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (m *MonoBankParser) GetBankName() types.BankName {
	return types.BankNameMono
}
