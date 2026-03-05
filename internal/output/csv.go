package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

// PrintResultsCSV writes lookup results as CSV with a header row.
func PrintResultsCSV(w io.Writer, results []LookupResult) {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	cw.Write([]string{"routing_number", "institution", "rtp", "fednow"})

	for _, r := range results {
		cw.Write([]string{
			r.RoutingNumber,
			r.Institution,
			fmt.Sprintf("%t", r.RTP),
			fmt.Sprintf("%t", r.FedNow),
		})
	}
}

// PrintAnalysisCSV is an alias for PrintResultsCSV.
func PrintAnalysisCSV(w io.Writer, results []LookupResult) {
	PrintResultsCSV(w, results)
}
