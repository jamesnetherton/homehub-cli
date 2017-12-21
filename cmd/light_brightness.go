package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewLightBrightnessCommand creates a new command to invoke the Hub LightBrightness function
func NewLightBrightnessCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightBrightness",
			Description: "Gets the Home Hub LED brightness percentage value",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().LightBrightness()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
