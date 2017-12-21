package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewWiFiSSIDCommand creates a new command to invoke the Hub WiFiSSID function
func NewWiFiSSIDCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "WiFiSSID",
			Description: "Gets the Home Hub WiFI SSID",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().WiFiSSID()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
