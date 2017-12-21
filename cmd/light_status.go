package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewLightStatusCommand creates a new command to invoke the Hub LightStatus function
func NewLightStatusCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightStatus",
			Description: "Gets the status of the Home Hub LED",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().LightStatus()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
