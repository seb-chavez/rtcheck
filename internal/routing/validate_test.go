package routing

import "testing"

func TestNormalizeRoutingNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  021000021  ", "021000021"},
		{"021-000-021", "021000021"},
		{"21000021", "021000021"},
		{"abc", "abc"},
		{"123", "000000123"},
		{"", ""},
	}
	for _, tt := range tests {
		got := Normalize(tt.input)
		if got != tt.expected {
			t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsValid(t *testing.T) {
	valid := []string{"021000021", "031100209", "071000013"}
	for _, rtn := range valid {
		if !IsValid(rtn) {
			t.Errorf("IsValid(%q) = false, want true", rtn)
		}
	}

	invalid := []string{"123456789", "12345", "abcdefghi", "000000000", ""}
	for _, rtn := range invalid {
		if IsValid(rtn) {
			t.Errorf("IsValid(%q) = true, want false", rtn)
		}
	}
}
