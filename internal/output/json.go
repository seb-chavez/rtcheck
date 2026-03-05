package output

import (
	"encoding/json"
	"io"
)

// PrintLookupJSON writes a single LookupResult as indented JSON.
func PrintLookupJSON(w io.Writer, r LookupResult) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(r)
}

// PrintAnalysisJSON writes an analysis summary and results as indented JSON.
func PrintAnalysisJSON(w io.Writer, summary AnalysisSummary, results []LookupResult) {
	output := struct {
		Summary AnalysisSummary `json:"summary"`
		Results []LookupResult  `json:"results"`
	}{
		Summary: summary,
		Results: results,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(output)
}

// PrintDirectoryJSON writes an array of LookupResults as indented JSON.
func PrintDirectoryJSON(w io.Writer, results []LookupResult) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(results)
}
