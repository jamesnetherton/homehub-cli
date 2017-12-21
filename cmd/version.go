package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewVersionCommand creates a new command to invoke the Hub Version function
func NewVersionCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "Version",
			Description: "Gets the Home Hub version",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().Version()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
