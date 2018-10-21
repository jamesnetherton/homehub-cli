package cmd

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestEventLogCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)

	eventLog := &homehub.EventLog{}
	var logEntries []homehub.EventLogEntry

	for i := 1; i <= 5; i++ {
		eventLogEntry := &homehub.EventLogEntry{
			Category:  "INFO",
			Message:   fmt.Sprintf("Log entry %d", i),
			Timestamp: "2018-01-01 00:00:00",
			Type:      fmt.Sprintf("Type %d", i),
		}
		logEntries = append(logEntries, *eventLogEntry)
	}

	eventLog.Entries = logEntries

	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().EventLog().Return(eventLog, nil)

	AssertCommandOutput(t, NewEventLogCommand(NewLoginCommand()))
}
