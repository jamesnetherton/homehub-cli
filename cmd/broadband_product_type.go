package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewBroadbandProductTypeCommand creates a new command to invoke the Hub BroadbandProductType function
func NewBroadbandProductTypeCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "BroadbandProductType",
			Description: "Gets the BT broadband product type",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().BroadbandProductType()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
