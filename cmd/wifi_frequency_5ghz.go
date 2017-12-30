package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewWiFiFrquency5GhzCommand creates a new command to invoke the Hub WiFiFrquency5GhzCommand function
func NewWiFiFrquency5GhzCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "WiFiFrequency5Ghz",
			Description: "Gets information relating to the Home Hub 5GHz wireless frequency",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().WiFiFrequency5Ghz()) },
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					dataPattern := "%-25s%-5d\n%-25s%-25s\n%-25s%-25t\n%-25s%-25s\n%-25s%-15d\n%-25s%-25s\n%-25s%-25s\n%-25s%-25s\n"
					wiFiFrequency := context.GetResult().(*homehub.WiFiFrequency)

					fmt.Print("\n")
					fmt.Printf(dataPattern,
						"ID",
						wiFiFrequency.UID,
						"Alias",
						wiFiFrequency.Alias,
						"Enabled",
						wiFiFrequency.Enable,
						"Status",
						wiFiFrequency.Status,
						"Channel",
						wiFiFrequency.Channel,
						"Available Channels",
						wiFiFrequency.AvailableChannels,
						"Operating Standards",
						wiFiFrequency.OperatingStandards,
						"Supported Standards",
						wiFiFrequency.SupportedStandards,
					)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
