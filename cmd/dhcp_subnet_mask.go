package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDhcpSubnetMaskCommand creates a new command to invoke the Hub DhcpSubnetMask function
func NewDhcpSubnetMaskCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DhcpSubnetMask",
			Description: "Gets the Home Hub DHCP subnet mask",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DhcpSubnetMask()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
