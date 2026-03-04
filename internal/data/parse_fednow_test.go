package data

import "testing"

func TestParseFedNowData(t *testing.T) {
	input := "021000021\r\n031100209\r\n071000013\r\n"
	rtns := ParseFedNowData([]byte(input))
	if len(rtns) != 3 {
		t.Fatalf("ParseFedNowData returned %d entries, want 3", len(rtns))
	}
	expected := []string{"021000021", "031100209", "071000013"}
	for i, rtn := range rtns {
		if rtn != expected[i] {
			t.Errorf("rtns[%d] = %q, want %q", i, rtn, expected[i])
		}
	}
}
