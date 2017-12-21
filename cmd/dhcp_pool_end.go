package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDhcpPoolEndCommand creates a new command to invoke the Hub DhcpPoolEnd function
func NewDhcpPoolEndCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DhcpPoolEnd",
			Description: "Gets the Home Hub IPV4 DHCP pool end address",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DhcpPoolEnd()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
