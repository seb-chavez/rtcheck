package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

// PrintResultsCSV writes lookup results as CSV with a header row.
// It returns an error if writing to the underlying writer fails.
func PrintResultsCSV(w io.Writer, results []LookupResult) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"routing_number", "institution", "rtp", "fednow"}); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	for _, r := range results {
		if err := cw.Write([]string{
			r.RoutingNumber,
			r.Institution,
			fmt.Sprintf("%t", r.RTP),
			fmt.Sprintf("%t", r.FedNow),
		}); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		return fmt.Errorf("flushing CSV writer: %w", err)
	}

	return nil
}
