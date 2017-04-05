package homehub

import (
	"strconv"
	"strings"
)

// BandwidthLog defines a Home Hub device bandwidth log history
type BandwidthLog struct {
	Entries []BandwidthLogEntry
}

// BandwidthLogEntry defines a Home Hub bandwidth log entry
type BandwidthLogEntry struct {
	MACAddress        string
	Date              string
	DownloadMegabytes int
	UploadMegabytes   int
}

func parseBandwidthLogEntry(logEntry string) *BandwidthLogEntry {
	entry := strings.Split(logEntry, ",")

	if len(entry) >= 5 {
		var messageParts []string
		for i := 0; i < len(entry); i++ {
			messageParts = append(messageParts, entry[i])
		}

		downloaded, downloadedErr := strconv.Atoi(entry[3])
		if downloadedErr != nil {
			downloaded = 0
		}

		uploaded, uploadedErr := strconv.Atoi(entry[4])
		if uploadedErr != nil {
			uploaded = 0
		}

		return &BandwidthLogEntry{
			entry[1],
			entry[2],
			downloaded,
			uploaded,
		}
	}
	return nil
}
