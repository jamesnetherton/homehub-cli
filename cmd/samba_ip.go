package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewSambaIPCommand creates a new command to invoke the Hub SambaIP function
func NewSambaIPCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "SambaIP",
			Description: "Gets the Home Hub samba IP address",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().SambaIP()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
