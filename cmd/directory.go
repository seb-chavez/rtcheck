package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/seb-chavez/rtcheck/internal/cache"
	"github.com/seb-chavez/rtcheck/internal/data"
	"github.com/seb-chavez/rtcheck/internal/output"
	"github.com/spf13/cobra"
)

var (
	searchTerm string
	network    string
)

var directoryCmd = &cobra.Command{
	Use:   "directory",
	Short: "Browse real-time payment participants",
	Long:  "Browse and search all routing numbers participating in RTP and/or FedNow.",
	RunE:  runDirectory,
}

func init() {
	rootCmd.AddCommand(directoryCmd)
	directoryCmd.Flags().StringVar(&formatOut, "format", "table", "Output format: table, json, csv")
	directoryCmd.Flags().StringVar(&searchTerm, "search", "", "Filter by institution name or routing number prefix")
	directoryCmd.Flags().StringVar(&network, "network", "", "Filter by network: rtp, fednow, both")
}

func runDirectory(cmd *cobra.Command, args []string) error {
	dir := cacheDir
	if dir == "" {
		dir = cache.DefaultDir()
	}
	c := cache.New(dir, 24*time.Hour)

	store, err := data.LoadStore(c, refresh)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}

	all := store.All()

	var filtered []output.LookupResult
	for _, inst := range all {
		// Network filter
		if network != "" {
			switch network {
			case "rtp":
				if !inst.RTP {
					continue
				}
			case "fednow":
				if !inst.FedNow {
					continue
				}
			case "both":
				if !inst.RTP || !inst.FedNow {
					continue
				}
			}
		}

		// Search filter
		if searchTerm != "" {
			nameMatch := strings.Contains(strings.ToLower(inst.Name), strings.ToLower(searchTerm))
			rtnMatch := strings.HasPrefix(inst.RoutingNumber, searchTerm)
			if !nameMatch && !rtnMatch {
				continue
			}
		}

		filtered = append(filtered, output.LookupResult{
			RoutingNumber: inst.RoutingNumber,
			Institution:   inst.Name,
			RTP:           inst.RTP,
			FedNow:        inst.FedNow,
		})
	}

	format := output.ParseFormat(formatOut)

	switch format {
	case output.FormatJSON:
		output.PrintDirectoryJSON(os.Stdout, filtered)
		return nil
	case output.FormatCSV:
		output.PrintResultsCSV(os.Stdout, filtered)
		return nil
	}

	// Table format: paginated interactive display
	pageSize := 50
	page := 1
	totalPages := (len(filtered) + pageSize - 1) / pageSize

	scanner := bufio.NewScanner(os.Stdin)
	for {
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > len(filtered) {
			end = len(filtered)
		}

		output.PrintDirectoryTable(os.Stdout, filtered[start:end], page, pageSize, len(filtered))

		if totalPages <= 1 {
			break
		}

		fmt.Printf("  [n]ext [p]rev [q]uit: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		switch input {
		case "n", "next":
			if page < totalPages {
				page++
			}
		case "p", "prev":
			if page > 1 {
				page--
			}
		case "q", "quit", "":
			return nil
		}
	}

	return nil
}
