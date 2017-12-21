package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDataReceivedCommand creates a new command to invoke the Hub DataReceived function
func NewDataReceivedCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DataReceived",
			Description: "Gets the number of bytes receieved since the Home Hub was last rebooted",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DataReceived()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
