package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewLightEnableCommand creates a new command to invoke the Hub LightEnable function
func NewLightEnableCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightEnable",
			Description: "Toggles the Home Hub LED on or off",
			Exec: func(context *CommandContext) {
				enable, _ := context.GetBooleanArg(0)
				context.SetResult(nil, service.GetHub().LightEnable(enable))
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
