package cmd

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestBandwidthMonitorCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)

	bandwidthLog := &homehub.BandwidthLog{}
	var logEntries []homehub.BandwidthLogEntry

	for i := 1; i <= 5; i++ {
		logEntry := &homehub.BandwidthLogEntry{
			MACAddress:        fmt.Sprintf("AA:BB:CC:DD:EE:FF:0%d", i),
			Date:              fmt.Sprintf("2018-01-0%d", i),
			DownloadMegabytes: i * 1000,
			UploadMegabytes:   i * 2000,
		}
		logEntries = append(logEntries, *logEntry)
	}

	bandwidthLog.Entries = logEntries

	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().BandwidthMonitor().Return(bandwidthLog, nil)

	AssertCommandOutput(t, NewBandwidthMonitorCommand(NewLoginCommand()))
}
