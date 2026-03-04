package routing

import (
	"regexp"
	"strings"
)

var dashSpaceRegex = regexp.MustCompile(`[-\s]`)
var digitsOnlyRegex = regexp.MustCompile(`^\d{1,9}$`)
var nineDigitsRegex = regexp.MustCompile(`^\d{9}$`)

func Normalize(input string) string {
	cleaned := strings.TrimSpace(input)
	cleaned = dashSpaceRegex.ReplaceAllString(cleaned, "")
	if digitsOnlyRegex.MatchString(cleaned) {
		for len(cleaned) < 9 {
			cleaned = "0" + cleaned
		}
	}
	return cleaned
}

func IsValid(rtn string) bool {
	if !nineDigitsRegex.MatchString(rtn) {
		return false
	}
	if rtn == "000000000" {
		return false
	}
	d := make([]int, 9)
	for i, ch := range rtn {
		d[i] = int(ch - '0')
	}
	checksum := 3*(d[0]+d[3]+d[6]) + 7*(d[1]+d[4]+d[7]) + (d[2]+d[5]+d[8])
	return checksum%10 == 0
}
