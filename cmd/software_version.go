package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewSoftwareVersionCommand creates a new command to invoke the Hub SoftwareVersion function
func NewSoftwareVersionCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "SoftwareVersion",
			Description: "Gets the Home Hub software version",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().SoftwareVersion()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
