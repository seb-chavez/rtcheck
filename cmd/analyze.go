package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/seb-chavez/rtcheck/internal/cache"
	"github.com/seb-chavez/rtcheck/internal/data"
	"github.com/seb-chavez/rtcheck/internal/fileparse"
	"github.com/seb-chavez/rtcheck/internal/output"
	"github.com/seb-chavez/rtcheck/internal/routing"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	noSummary  bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <file>",
	Short: "Analyze a file of routing numbers",
	Long:  "Bulk-check a CSV or Excel file of routing numbers against RTP and FedNow networks.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVar(&formatOut, "format", "table", "Output format: table, json, csv")
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write detailed results to file")
	analyzeCmd.Flags().BoolVar(&noSummary, "no-summary", false, "Skip summary, output detail rows only")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	filename := args[0]
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	parsed, err := fileparse.Parse(fileData, filepath.Base(filename))
	if err != nil {
		return err
	}

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, rtn := range parsed.RoutingNumbers {
		normalized := routing.Normalize(rtn)
		if !seen[normalized] && routing.IsValid(normalized) {
			seen[normalized] = true
			unique = append(unique, normalized)
		}
	}

	if len(unique) == 0 {
		return fmt.Errorf("no valid routing numbers found in %s", filename)
	}

	dir := cacheDir
	if dir == "" {
		dir = cache.DefaultDir()
	}
	c := cache.New(dir, 24*time.Hour)

	store, err := data.LoadStore(c, refresh)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}

	var results []output.LookupResult
	rtpCount, fnCount, bothCount := 0, 0, 0
	for _, rtn := range unique {
		inst := store.Lookup(rtn)
		results = append(results, output.LookupResult{
			RoutingNumber: inst.RoutingNumber,
			Institution:   inst.Name,
			RTP:           inst.RTP,
			FedNow:        inst.FedNow,
		})
		if inst.RTP {
			rtpCount++
		}
		if inst.FedNow {
			fnCount++
		}
		if inst.RTP && inst.FedNow {
			bothCount++
		}
	}

	total := len(unique)
	neitherCount := total - (rtpCount + fnCount - bothCount)

	pct := func(count int) float64 {
		if total == 0 {
			return 0
		}
		return math.Round(float64(count)/float64(total)*1000) / 10
	}

	summary := output.AnalysisSummary{
		File:           filepath.Base(filename),
		TotalUnique:    total,
		RTPCount:       rtpCount,
		FedNowCount:    fnCount,
		BothCount:      bothCount,
		NeitherCount:   neitherCount,
		RTPPercent:     pct(rtpCount),
		FedNowPercent:  pct(fnCount),
		BothPercent:    pct(bothCount),
		NeitherPercent: pct(neitherCount),
	}

	format := output.ParseFormat(formatOut)

	switch format {
	case output.FormatJSON:
		if err := output.PrintAnalysisJSON(os.Stdout, summary, results); err != nil {
			return fmt.Errorf("writing JSON output: %w", err)
		}
	case output.FormatCSV:
		if err := output.PrintAnalysisCSV(os.Stdout, results); err != nil {
			return fmt.Errorf("writing CSV output: %w", err)
		}
	default:
		if !noSummary {
			output.PrintAnalysisSummaryTable(os.Stdout, summary)
		}
		if noSummary {
			output.PrintResultsTable(os.Stdout, results)
		}
	}

	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		if err := output.PrintResultsCSV(f, results); err != nil {
			return fmt.Errorf("writing CSV to file: %w", err)
		}
		if format == output.FormatTable && !noSummary {
			fmt.Fprintf(os.Stdout, "  Detailed results written to: %s\n\n", outputFile)
		}
	}

	return nil
}
