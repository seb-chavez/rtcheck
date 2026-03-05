package output

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

const (
	checkMark = "\u2713"
	crossMark = "\u2717"
)

func boolToStatus(b bool) string {
	if b {
		return checkMark
	}
	return crossMark
}

func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// PrintLookupTable prints a formatted block showing routing number details.
func PrintLookupTable(w io.Writer, r LookupResult) {
	realTime := r.RTP || r.FedNow

	table := tablewriter.NewWriter(w)
	table.Header("FIELD", "VALUE")
	table.Append("Routing Number", r.RoutingNumber)
	table.Append("Institution", r.Institution)
	table.Append("RTP", boolToYesNo(r.RTP))
	table.Append("FedNow", boolToYesNo(r.FedNow))
	table.Append("Real-Time Capable", boolToYesNo(realTime))
	table.Render()
}

// PrintAnalysisSummaryTable prints a summary table with Network/Count/Coverage columns.
func PrintAnalysisSummaryTable(w io.Writer, s AnalysisSummary) {
	fmt.Fprintf(w, "\nFile: %s\n", s.File)
	fmt.Fprintf(w, "Total unique routing numbers: %d\n\n", s.TotalUnique)

	table := tablewriter.NewWriter(w)
	table.Header("NETWORK", "COUNT", "COVERAGE")
	table.Append("RTP", fmt.Sprintf("%d", s.RTPCount), fmt.Sprintf("%.1f%%", s.RTPPercent))
	table.Append("FedNow", fmt.Sprintf("%d", s.FedNowCount), fmt.Sprintf("%.1f%%", s.FedNowPercent))
	table.Append("Both", fmt.Sprintf("%d", s.BothCount), fmt.Sprintf("%.1f%%", s.BothPercent))
	table.Append("Neither", fmt.Sprintf("%d", s.NeitherCount), fmt.Sprintf("%.1f%%", s.NeitherPercent))
	table.Render()
}

// PrintDirectoryTable prints a paginated table of routing numbers with checkmarks for network support.
func PrintDirectoryTable(w io.Writer, results []LookupResult, page, pageSize, total int) {
	table := tablewriter.NewWriter(w)
	table.Header("ROUTING NUMBER", "INSTITUTION", "RTP", "FEDNOW")

	for _, r := range results {
		table.Append(r.RoutingNumber, r.Institution, boolToStatus(r.RTP), boolToStatus(r.FedNow))
	}

	table.Render()

	totalPages := (total + pageSize - 1) / pageSize
	fmt.Fprintf(w, "\nPage %d of %d (showing %d of %d results)\n", page, totalPages, len(results), total)
}

// PrintResultsTable prints a detailed results table to stdout.
func PrintResultsTable(results []LookupResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ROUTING NUMBER", "INSTITUTION", "RTP", "FEDNOW")

	for _, r := range results {
		table.Append(r.RoutingNumber, r.Institution, boolToStatus(r.RTP), boolToStatus(r.FedNow))
	}

	table.Render()
}
