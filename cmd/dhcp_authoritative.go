package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDhcpAuthoritativeCommand creates a new command to invoke the Hub DhcpAuthoritative function
func NewDhcpAuthoritativeCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DhcpAuthoritative",
			Description: "Gets details about whether the Home Hub is the authoritative DHCP server",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DhcpAuthoritative()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
