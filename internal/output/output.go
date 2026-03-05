package output

// Format represents the output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
)

// ParseFormat converts a string to a Format, defaulting to FormatTable.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "csv":
		return FormatCSV
	default:
		return FormatTable
	}
}

// LookupResult holds the real-time payment network support for a routing number.
type LookupResult struct {
	RoutingNumber string `json:"routing_number"`
	Institution   string `json:"institution,omitempty"`
	RTP           bool   `json:"rtp"`
	FedNow        bool   `json:"fednow"`
}

// AnalysisSummary holds aggregate statistics from analyzing a file of routing numbers.
type AnalysisSummary struct {
	File           string  `json:"file"`
	TotalUnique    int     `json:"total_unique"`
	RTPCount       int     `json:"rtp_count"`
	FedNowCount    int     `json:"fednow_count"`
	BothCount      int     `json:"both_count"`
	NeitherCount   int     `json:"neither_count"`
	RTPPercent     float64 `json:"rtp_percent"`
	FedNowPercent  float64 `json:"fednow_percent"`
	BothPercent    float64 `json:"both_percent"`
	NeitherPercent float64 `json:"neither_percent"`
}
