package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewWiFiFrequency5GhzChannelSetCommand creates a new command to invoke the Hub WiFiFrequency5GhzChannelSet function
func NewWiFiFrequency5GhzChannelSetCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "WiFiFrequency5GhzChannelSet",
			Description: "Sets the Home Hub 5GHz wireless frequency channel",
			ArgNames:    []string{"channel"},
			ArgTypes:    []string{"int"},
			Exec: func(context *CommandContext) {
				channel, err := context.GetIntArg(0)
				if err != nil || channel <= 0 {
					parseErr := errors.New("Channel must be a positive numeric value")
					context.SetResult(nil, parseErr)
					return
				}

				context.SetResult(nil, service.GetHub().WiFiFrequency5GhzChannelSet(channel))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					fmt.Printf("WiFi 5GHz channel updated successfully\n")
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
