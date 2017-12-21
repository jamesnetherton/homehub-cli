package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewSambaHostCommand creates a new command to invoke the Hub SambaHost function
func NewSambaHostCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "SambaHost",
			Description: "Gets the Home Hub samba host name",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().SambaHost()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
