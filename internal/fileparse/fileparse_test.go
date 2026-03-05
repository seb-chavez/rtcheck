package fileparse

import (
	"os"
	"testing"
)

func TestParseCSVWithHeader(t *testing.T) {
	data, err := os.ReadFile("../../testdata/sample.csv")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	result, err := Parse(data, "sample.csv")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(result.RoutingNumbers) != 10 {
		t.Errorf("got %d routing numbers, want 10", len(result.RoutingNumbers))
	}
	if result.AccountCounts == nil {
		t.Error("expected AccountCounts to be non-nil")
	}
	if result.AccountCounts["021000021"] != 1500 {
		t.Errorf("AccountCounts[021000021] = %d, want 1500", result.AccountCounts["021000021"])
	}
}

func TestParsePlainText(t *testing.T) {
	data, err := os.ReadFile("../../testdata/plain.txt")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	result, err := Parse(data, "plain.txt")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(result.RoutingNumbers) != 3 {
		t.Errorf("got %d routing numbers, want 3", len(result.RoutingNumbers))
	}
}

func TestParseAutoDetectColumn(t *testing.T) {
	csv := "id,bank_code,amount\n1,021000021,500\n2,031100209,300\n3,071000013,100\n"
	result, err := Parse([]byte(csv), "test.csv")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(result.RoutingNumbers) != 3 {
		t.Errorf("got %d routing numbers, want 3", len(result.RoutingNumbers))
	}
}

func TestParseHeaderVariants(t *testing.T) {
	variants := []string{"routing_number", "RTN", "ABA", "routing number", "bank_routing", "transit_number", "r/t"}
	for _, header := range variants {
		csv := header + "\n021000021\n031100209\n"
		result, err := Parse([]byte(csv), "test.csv")
		if err != nil {
			t.Errorf("header %q: Parse failed: %v", header, err)
			continue
		}
		if len(result.RoutingNumbers) != 2 {
			t.Errorf("header %q: got %d routing numbers, want 2", header, len(result.RoutingNumbers))
		}
	}
}
