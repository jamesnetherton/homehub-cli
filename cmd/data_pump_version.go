package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewDataPumpVersionCommand creates a new command to invoke the Hub DataPumpVersion function
func NewDataPumpVersionCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DataPumpVersion",
			Description: "Gets details related to the DSL line firmware version",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().DataPumpVersion()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
