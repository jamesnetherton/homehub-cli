package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewMaintenanceFirmwareVersionCommand creates a new command to invoke the Hub MaintenanceFirmwareVersion function
func NewMaintenanceFirmwareVersionCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "MaintenanceFirmwareVersion",
			Description: "Gets the Home Hub maintenance firmware version",
			Exec: func(context *CommandContext) {
				context.SetResult(service.GetHub().MaintenanceFirmwareVersion())
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
