package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewBandwidthMonitorCommand creates a new command to invoke the Hub BandwidthMonitor function
func NewBandwidthMonitorCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "BandwidthMonitor",
			Description: "Displays bandwidth statistics for devices that have connected to the Home Hub",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().BandwidthMonitor()) },
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					headerPattern := "%-30s%-20s%-25s%-7s\n"
					dataPattern := "%-30s%-20s%-25d%-7d\n"
					bandwidthLog := context.GetResult().(*homehub.BandwidthLog)
					bandwidthLogEntries := bandwidthLog.Entries

					fmt.Print("\n")
					fmt.Printf(headerPattern, "------------------", "----------", "----------", "--------")
					fmt.Printf(headerPattern, "   MAC Address    ", "   Date   ", "Downloaded", "Uploaded")
					fmt.Printf(headerPattern, "------------------", "----------", "----------", "--------")

					for i := 0; i < len(bandwidthLogEntries); i++ {
						fmt.Printf(dataPattern,
							bandwidthLogEntries[i].MACAddress,
							bandwidthLogEntries[i].Date,
							bandwidthLogEntries[i].DownloadMegabytes,
							bandwidthLogEntries[i].UploadMegabytes,
						)
					}

					fmt.Print("\n")
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
