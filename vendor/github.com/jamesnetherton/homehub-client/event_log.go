package homehub

import (
	"fmt"
	"strings"
)

// EventLog defines a Home Hub event log history
type EventLog struct {
	Entries []EventLogEntry
}

// EventLogEntry defines a Home Hub event log entry
type EventLogEntry struct {
	Timestamp string
	Type      string
	Category  string
	Message   string
}

func parseEventLogEntry(logEntry string) *EventLogEntry {
	entry := strings.Split(logEntry, " ")

	if len(entry) >= 6 {
		var messageParts []string
		for i := 5; i < len(entry); i++ {
			messageParts = append(messageParts, entry[i])
		}

		return &EventLogEntry{
			fmt.Sprintf("%s %s", entry[0], entry[1]),
			entry[2],
			entry[3],
			strings.Join(messageParts, " "),
		}
	}
	return nil
}
