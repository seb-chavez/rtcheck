package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/seb-chavez/rtcheck/internal/cache"
	"github.com/seb-chavez/rtcheck/internal/data"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	refresh   bool
	cacheDir  string
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
}

func loadStore() (*data.Store, error) {
	dir := cacheDir
	if dir == "" {
		dir = cache.DefaultDir()
	}
	c := cache.New(dir, 24*time.Hour)
	return data.LoadStore(c, refresh)
}
