package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewLocalTimeCommand creates a new command to invoke the Hub LocalTime function
func NewLocalTimeCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LocalTime",
			Description: "Gets local time from the Home Hub",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().LocalTime()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
