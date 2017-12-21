package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewSerialNumberCommand creates a new command to invoke the Hub SerialNumber function
func NewSerialNumberCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "SerialNumber",
			Description: "Gets the Home Hub serial number",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().SerialNumber()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
