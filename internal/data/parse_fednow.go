package data

// ParseFedNowData extracts 9-digit routing numbers from raw FedNow participant data.
// The format is identical to RTP data, so this delegates to ParseRTPData.
func ParseFedNowData(raw []byte) []string {
	return ParseRTPData(raw)
}
