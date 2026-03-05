package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/moov-io/fed"
	"github.com/seb-chavez/rtcheck/internal/cache"
)

const (
	rtpPageURL    = "https://www.theclearinghouse.org/payment-systems/rtp/rtn"
	fednowPageURL = "https://www.frbservices.org/financial-services/fednow/organizations/terms-of-use"

	rtpCacheFile    = "rtp_rtns.json"
	fednowCacheFile = "fednow_rtns.json"
)

// cachedData is the JSON structure stored in the cache.
type cachedData struct {
	FetchedAt      string   `json:"fetched_at"`
	SourceURL      string   `json:"source_url"`
	RoutingNumbers []string `json:"routing_numbers"`
}

// LoadStore builds a Store populated with RTP, FedNow, and institution name data.
// It reads from the cache when available, fetching fresh data from official sources
// when the cache is stale or forceRefresh is true.
func LoadStore(c *cache.Cache, forceRefresh bool) (*Store, error) {
	store := NewStore()

	// Load RTP data
	rtpRTNs, err := loadOrFetch(c, rtpCacheFile, rtpPageURL, "TXT.txt", ParseRTPData, forceRefresh)
	if err != nil {
		return nil, fmt.Errorf("loading RTP data: %w", err)
	}
	store.LoadRTP(rtpRTNs)

	// Load FedNow data
	fnRTNs, err := loadOrFetch(c, fednowCacheFile, fednowPageURL, ".txt", ParseFedNowData, forceRefresh)
	if err != nil {
		return nil, fmt.Errorf("loading FedNow data: %w", err)
	}
	store.LoadFedNow(fnRTNs)

	// Load institution names from moov-io/fed (graceful degradation)
	names := loadFedACHNames()
	if names != nil {
		store.LoadNames(names)
	}

	return store, nil
}

// loadOrFetch attempts to read cached data. If the cache is stale, missing,
// or forceRefresh is true, it scrapes the page for a download link, fetches the
// file, parses it, and caches the result.
func loadOrFetch(
	c *cache.Cache,
	cacheFile string,
	pageURL string,
	linkSubstr string,
	parser func([]byte) []string,
	forceRefresh bool,
) ([]string, error) {
	// Try cache first
	if !forceRefresh {
		data, err := c.Read(cacheFile)
		if err == nil {
			var cd cachedData
			if json.Unmarshal(data, &cd) == nil && len(cd.RoutingNumbers) > 0 {
				return cd.RoutingNumbers, nil
			}
		}
		// cache miss or expired — fall through to fetch
	}

	// Scrape the page for a download link
	pageHTML, err := httpGet(pageURL)
	if err != nil {
		return nil, fmt.Errorf("fetching page %s: %w", pageURL, err)
	}

	downloadURL := extractLink(string(pageHTML), linkSubstr)
	if downloadURL == "" {
		return nil, fmt.Errorf("no download link containing %q found on %s", linkSubstr, pageURL)
	}

	// Make relative URLs absolute
	if strings.HasPrefix(downloadURL, "/") {
		// Extract base URL (scheme + host)
		parts := strings.SplitN(pageURL, "//", 2)
		if len(parts) == 2 {
			hostEnd := strings.Index(parts[1], "/")
			if hostEnd == -1 {
				hostEnd = len(parts[1])
			}
			downloadURL = parts[0] + "//" + parts[1][:hostEnd] + downloadURL
		}
	}

	// Download the file
	raw, err := httpGet(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("downloading %s: %w", downloadURL, err)
	}

	// Parse
	rtns := parser(raw)
	if len(rtns) == 0 {
		return nil, fmt.Errorf("parsed 0 routing numbers from %s", downloadURL)
	}

	// Cache the result
	cd := cachedData{
		FetchedAt:      time.Now().UTC().Format(time.RFC3339),
		SourceURL:      downloadURL,
		RoutingNumbers: rtns,
	}
	encoded, err := json.MarshalIndent(cd, "", "  ")
	if err == nil {
		if writeErr := c.Write(cacheFile, encoded); writeErr != nil {
			log.Printf("warning: failed to write cache file %s: %v", cacheFile, writeErr)
		}
	}

	return rtns, nil
}

// extractLink finds the first href in the HTML whose value contains substr.
func extractLink(html, substr string) string {
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) >= 2 && strings.Contains(m[1], substr) {
			return m[1]
		}
	}
	return ""
}

// httpGet performs an HTTP GET and returns the response body.
func httpGet(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

// loadFedACHNames attempts to load institution names from the moov-io/fed
// FedACH dictionary. It tries several common file paths. If the data cannot
// be loaded, it logs a warning and returns nil (graceful degradation).
func loadFedACHNames() map[string]string {
	// Try paths in order of preference
	paths := fedACHPaths()

	for _, path := range paths {
		names, err := func() (map[string]string, error) {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			dict := fed.NewACHDictionary()
			if err := dict.Read(f); err != nil {
				return nil, fmt.Errorf("failed to parse FedACH dictionary at %s: %w", path, err)
			}

			names := make(map[string]string, len(dict.ACHParticipants))
			for _, p := range dict.ACHParticipants {
				names[p.RoutingNumber] = p.CustomerName
			}
			return names, nil
		}()
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
		if names != nil {
			return names
		}
	}

	log.Println("warning: FedACH dictionary not found; institution names will be unavailable")
	return nil
}

// fedACHPaths returns candidate paths for the FedACH data file.
func fedACHPaths() []string {
	var paths []string

	// Environment variable override
	if envPath := os.Getenv("FEDACH_DATA_PATH"); envPath != "" {
		paths = append(paths, envPath)
	}

	// Common locations
	home, _ := os.UserHomeDir()
	if home != "" {
		paths = append(paths,
			home+"/.rtcheck/data/FedACHdir.txt",
			home+"/.rtcheck/FedACHdir.txt",
		)
	}

	// moov-io/fed module cache (Go modules)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = home + "/go"
	}

	// Try to find the moov-io/fed data file in the module cache
	// We check for the specific version that was installed
	modBase := gopath + "/pkg/mod/github.com/moov-io"
	entries, err := os.ReadDir(modBase)
	if err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "fed@") {
				candidate := modBase + "/" + entry.Name() + "/data/FedACHdir.txt"
				paths = append(paths, candidate)
			}
		}
	}

	// Current directory
	paths = append(paths, "FedACHdir.txt", "data/FedACHdir.txt")

	return paths
}
