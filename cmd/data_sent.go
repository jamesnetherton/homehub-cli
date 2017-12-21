package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDataSentCommand creates a new command to invoke the Hub DataSent function
func NewDataSentCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DataSent",
			Description: "Gets the number of bytes sent since the Home Hub was last rebooted",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DataSent()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
