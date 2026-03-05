package fileparse

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/seb-chavez/rtcheck/internal/routing"
	"github.com/xuri/excelize/v2"
)

// ParseResult holds the extracted routing numbers and optional account counts.
type ParseResult struct {
	RoutingNumbers []string
	AccountCounts  map[string]int // nil if no account_count column found
}

// Regex patterns for detecting routing number column headers.
var routingHeaderPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^routing[_\s-]?n(um(ber)?)?$`),
	regexp.MustCompile(`(?i)^rtn$`),
	regexp.MustCompile(`(?i)^aba$`),
	regexp.MustCompile(`(?i)^aba[_\s-]?r(outing)?[_\s-]?n(um(ber)?)?$`),
	regexp.MustCompile(`(?i)^bank[_\s-]?routing$`),
	regexp.MustCompile(`(?i)^routing$`),
	regexp.MustCompile(`(?i)^transit[_\s-]?n(um(ber)?)?$`),
	regexp.MustCompile(`(?i)^r/t$`),
	regexp.MustCompile(`(?i)^rt$`),
	regexp.MustCompile(`(?i)^routing[_\s-]?transit$`),
}

var nineDigitRegex = regexp.MustCompile(`^\d{9}$`)

// Parse detects file type by extension and extracts routing numbers.
// Supported formats: .csv, .xlsx, .xls, and plain text.
func Parse(data []byte, filename string) (*ParseResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".xlsx", ".xls":
		return parseExcel(data)
	case ".csv":
		return parseCSV(data)
	default:
		return parseMaybeTextOrCSV(data)
	}
}

// parseMaybeTextOrCSV tries plain text first, then falls back to CSV parsing.
func parseMaybeTextOrCSV(data []byte) (*ParseResult, error) {
	lines := splitLines(data)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	// Check if the first non-empty line looks like a 9-digit number (plain text mode).
	firstLine := strings.TrimSpace(lines[0])
	normalized := routing.Normalize(firstLine)
	if nineDigitRegex.MatchString(normalized) {
		return parsePlainText(lines)
	}

	// Otherwise try CSV parsing.
	return parseCSV(data)
}

// parsePlainText treats each line as a routing number.
func parsePlainText(lines []string) (*ParseResult, error) {
	result := &ParseResult{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		normalized := routing.Normalize(line)
		if nineDigitRegex.MatchString(normalized) {
			result.RoutingNumbers = append(result.RoutingNumbers, normalized)
		}
	}
	if len(result.RoutingNumbers) == 0 {
		return nil, fmt.Errorf("no routing numbers found in plain text")
	}
	return result, nil
}

// parseCSV parses CSV data, detects routing number and account count columns.
func parseCSV(data []byte) (*ParseResult, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV has fewer than 2 rows")
	}

	header := records[0]
	dataRows := records[1:]

	routingCol := findRoutingColumn(header)
	accountCountCol := findAccountCountColumn(header)

	// If no header match, try auto-detection on data rows.
	if routingCol < 0 {
		routingCol = autoDetectRoutingColumn(records)
	}

	if routingCol < 0 {
		return nil, fmt.Errorf("could not identify routing number column")
	}

	return extractFromRows(dataRows, routingCol, accountCountCol)
}

// parseExcel reads the first sheet of an Excel file and extracts routing numbers.
func parseExcel(data []byte) (*ParseResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("opening Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading Excel rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel sheet has fewer than 2 rows")
	}

	header := rows[0]
	dataRows := rows[1:]

	routingCol := findRoutingColumn(header)
	accountCountCol := findAccountCountColumn(header)

	if routingCol < 0 {
		routingCol = autoDetectRoutingColumn(rows)
	}

	if routingCol < 0 {
		return nil, fmt.Errorf("could not identify routing number column")
	}

	return extractFromRows(dataRows, routingCol, accountCountCol)
}

// findRoutingColumn checks header cells against known routing number patterns.
// Returns the column index or -1 if no match.
func findRoutingColumn(header []string) int {
	for i, h := range header {
		h = strings.TrimSpace(h)
		for _, pat := range routingHeaderPatterns {
			if pat.MatchString(h) {
				return i
			}
		}
	}
	return -1
}

// findAccountCountColumn finds an "account_count" column (case-insensitive).
// Returns the column index or -1 if not found.
func findAccountCountColumn(header []string) int {
	for i, h := range header {
		normalized := strings.ToLower(strings.TrimSpace(h))
		if normalized == "account_count" {
			return i
		}
	}
	return -1
}

// autoDetectRoutingColumn finds the column where >50% of values are valid 9-digit numbers.
// It skips the first row (assumed header) and checks all data rows.
func autoDetectRoutingColumn(allRows [][]string) int {
	if len(allRows) < 2 {
		return -1
	}

	header := allRows[0]
	dataRows := allRows[1:]

	numCols := len(header)
	bestCol := -1
	bestRatio := 0.5 // must exceed 50%

	for col := 0; col < numCols; col++ {
		matches := 0
		total := 0
		for _, row := range dataRows {
			if col >= len(row) {
				continue
			}
			total++
			normalized := routing.Normalize(strings.TrimSpace(row[col]))
			if nineDigitRegex.MatchString(normalized) {
				matches++
			}
		}
		if total > 0 {
			ratio := float64(matches) / float64(total)
			if ratio > bestRatio {
				bestRatio = ratio
				bestCol = col
			}
		}
	}

	return bestCol
}

// extractFromRows builds a ParseResult from data rows given column indices.
func extractFromRows(dataRows [][]string, routingCol, accountCountCol int) (*ParseResult, error) {
	result := &ParseResult{}

	if accountCountCol >= 0 {
		result.AccountCounts = make(map[string]int)
	}

	for _, row := range dataRows {
		if routingCol >= len(row) {
			continue
		}
		rtn := routing.Normalize(strings.TrimSpace(row[routingCol]))
		if !nineDigitRegex.MatchString(rtn) {
			continue
		}
		result.RoutingNumbers = append(result.RoutingNumbers, rtn)

		if accountCountCol >= 0 && accountCountCol < len(row) {
			countStr := strings.TrimSpace(row[accountCountCol])
			count, err := strconv.Atoi(countStr)
			if err == nil {
				result.AccountCounts[rtn] = count
			}
		}
	}

	if len(result.RoutingNumbers) == 0 {
		return nil, fmt.Errorf("no routing numbers found")
	}

	return result, nil
}

// splitLines splits data into lines, handling both \n and \r\n.
func splitLines(data []byte) []string {
	s := string(data)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}
