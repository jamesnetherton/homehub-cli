package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewWiFiSecurityModeCommand creates a new command to invoke the Hub WiFiSecurityMode function
func NewWiFiSecurityModeCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "WiFiSecurityMode",
			Description: "Gets the Home Hub WiFI security mode",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().WiFiSecurityMode()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
