package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewLightEnableCommand creates a new command to invoke the Hub LightEnable function
func NewLightEnableCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "LightEnable",
			Description: "Toggles the Home Hub LED on or off",
			ArgNames:    []string{"enable"},
			ArgTypes:    []string{"bool"},
			Exec: func(context *CommandContext) {
				enable, err := context.GetBooleanArg(0)
				if err != nil {
					parseErr := errors.New("LightEnable must be either true or false")
					context.SetResult(nil, parseErr)
					return
				}

				context.SetResult(nil, service.GetHub().LightEnable(enable))
			},
			PostExec: func(context *CommandContext) {
				enabled, _ := context.GetBooleanArg(0)
				if !context.IsError() {
					status := "disabled"
					if enabled {
						status = "enabled"
					}
					fmt.Printf("Light successfully %s\n", status)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
