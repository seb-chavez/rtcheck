package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	refresh   bool
	cacheDir  string
	noColor   bool
	formatOut string
)

var rootCmd = &cobra.Command{
	Use:     "rtcheck",
	Short:   "Check routing numbers against real-time payment networks",
	Long:    "rtcheck checks ABA routing numbers against RTP (The Clearing House) and FedNow (Federal Reserve) participant lists.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&refresh, "refresh", false, "Force re-download of participant data")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "Override cache directory (default: ~/.rtcheck/data/)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}
