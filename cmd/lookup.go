package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/seb-chavez/rtcheck/internal/cache"
	"github.com/seb-chavez/rtcheck/internal/data"
	"github.com/seb-chavez/rtcheck/internal/output"
	"github.com/seb-chavez/rtcheck/internal/routing"
	"github.com/spf13/cobra"
)

var lookupCmd = &cobra.Command{
	Use:   "lookup <routing-number>",
	Short: "Look up a single routing number",
	Long:  "Check if a routing number participates in RTP and/or FedNow networks.",
	Args:  cobra.ExactArgs(1),
	RunE:  runLookup,
}

func init() {
	rootCmd.AddCommand(lookupCmd)
	lookupCmd.Flags().StringVar(&formatOut, "format", "table", "Output format: table, json")
}

func runLookup(cmd *cobra.Command, args []string) error {
	rtn := routing.Normalize(args[0])

	if !routing.IsValid(rtn) {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid ABA routing number\n", args[0])
		os.Exit(1)
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

	inst := store.Lookup(rtn)
	result := output.LookupResult{
		RoutingNumber: inst.RoutingNumber,
		Institution:   inst.Name,
		RTP:           inst.RTP,
		FedNow:        inst.FedNow,
	}

	switch output.ParseFormat(formatOut) {
	case output.FormatJSON:
		output.PrintLookupJSON(os.Stdout, result)
	default:
		output.PrintLookupTable(os.Stdout, result)
	}

	if !inst.RTP && !inst.FedNow {
		os.Exit(2)
	}

	return nil
}
