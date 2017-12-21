package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDownstreamSyncSpeedCommand creates a new command to invoke the Hub DownstreamSyncSpeed function
func NewDownstreamSyncSpeedCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DownstreamSyncSpeed",
			Description: "Gets the available speed at which the Home Hub can download data",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DownstreamSyncSpeed()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
