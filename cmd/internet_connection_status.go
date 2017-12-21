package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewInternetConnectionStatusCommand creates a new command to invoke the Hub InternetConnectionStatus function
func NewInternetConnectionStatusCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "InternetConnectionStatus",
			Description: "Gets the status of the Home Hub internet connection",
			Exec: func(context *CommandContext) {
				context.SetResult(service.GetHub().InternetConnectionStatus())
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
