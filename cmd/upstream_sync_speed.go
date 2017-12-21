package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewUpstreamSyncSpeedCommand creates a new command to invoke the Hub UpstreamSyncSpeed function
func NewUpstreamSyncSpeedCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "UpstreamSyncSpeed",
			Description: "Gets the Home Hub upload speed",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().UpstreamSyncSpeed()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
