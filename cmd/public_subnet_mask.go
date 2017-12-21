package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewPublicSubnetMaskCommand creates a new command to invoke the Hub PublicSubnetMask function
func NewPublicSubnetMaskCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "PublicSubnetMask",
			Description: "Gets the Home Hub public IP subnet mask",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().PublicSubnetMask()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
