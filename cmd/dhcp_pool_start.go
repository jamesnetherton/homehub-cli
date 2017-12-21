package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDhcpPoolStartCommand creates a new command to invoke the Hub DhcpPoolStart function
func NewDhcpPoolStartCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DhcpPoolStart",
			Description: "Gets the Home Hub IPV4 DHCP pool start address",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DhcpPoolStart()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
