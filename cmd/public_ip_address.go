package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewPublicIPAddressCommand creates a new command to invoke the Hub PublicIPAddress function
func NewPublicIPAddressCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "PublicIPAddress",
			Description: "Gets the Home Hub public IP address",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().PublicIPAddress()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
