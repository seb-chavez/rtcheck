package data

import "strings"

func ParseFedNowData(raw []byte) []string {
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	var rtns []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if nineDigits.MatchString(trimmed) {
			rtns = append(rtns, trimmed)
		}
	}
	return rtns
}
