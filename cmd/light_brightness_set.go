package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewLightBrightnessSetCommand creates a new command to invoke the Hub LightBrightnessSet function
func NewLightBrightnessSetCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightBrightnessSet",
			Description: "Sets the Home Hub LED brightness percentage value",
			Exec: func(context *CommandContext) {
				brightness, _ := context.GetIntArg(0)
				context.SetResult(nil, service.GetHub().LightBrightnessSet(brightness))
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
