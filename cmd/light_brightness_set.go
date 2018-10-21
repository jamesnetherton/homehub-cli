package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewLightBrightnessSetCommand creates a new command to invoke the Hub LightBrightnessSet function
func NewLightBrightnessSetCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightBrightnessSet",
			Description: "Sets the Home Hub LED brightness percentage value",
			ArgNames:    []string{"brightness"},
			ArgTypes:    []string{"int"},
			Exec: func(context *CommandContext) {
				brightness, err := context.GetIntArg(0)
				if err != nil || brightness <= 0 || brightness > 100 {
					parseErr := errors.New("Brightness must be a numeric value between 0 and 100")
					context.SetResult(nil, parseErr)
					return
				}

				context.SetResult(nil, service.GetHub().LightBrightnessSet(brightness))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					fmt.Printf("Light brightness successfully updated\n")
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
