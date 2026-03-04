package data

import "testing"

func TestParseRTPData(t *testing.T) {
	input := "011000138\r011000206\r011000390\r"
	rtns := ParseRTPData([]byte(input))
	if len(rtns) != 3 {
		t.Fatalf("ParseRTPData returned %d entries, want 3", len(rtns))
	}
	expected := []string{"011000138", "011000206", "011000390"}
	for i, rtn := range rtns {
		if rtn != expected[i] {
			t.Errorf("rtns[%d] = %q, want %q", i, rtn, expected[i])
		}
	}
}

func TestParseRTPDataWithNewlines(t *testing.T) {
	input := "011000138\r\n011000206\n011000390\n"
	rtns := ParseRTPData([]byte(input))
	if len(rtns) != 3 {
		t.Fatalf("ParseRTPData returned %d entries, want 3", len(rtns))
	}
}

func TestParseRTPDataSkipsInvalid(t *testing.T) {
	input := "011000138\r\nnotanumber\r\n011000206\r\n\r\n"
	rtns := ParseRTPData([]byte(input))
	if len(rtns) != 2 {
		t.Fatalf("ParseRTPData returned %d entries, want 2", len(rtns))
	}
}
