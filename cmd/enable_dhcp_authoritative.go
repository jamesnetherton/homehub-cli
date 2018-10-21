package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewEnableDhcpAuthoritativeCommand creates a new command to invoke the Hub EnableDhcpAuthoritative function
func NewEnableDhcpAuthoritativeCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "EnableDhcpAuthoritative",
			Description: "Toggles whether the Home Hub is the authoritative DHCP server",
			ArgNames:    []string{"enable"},
			ArgTypes:    []string{"bool"},
			Exec: func(context *CommandContext) {
				enable, err := context.GetBooleanArg(0)
				if err != nil {
					parseErr := errors.New("Enable flag must be either true or false")
					context.SetResult(nil, parseErr)
					return
				}
				context.SetResult(nil, service.GetHub().EnableDhcpAuthoritative(enable))
			},
			PostExec: func(context *CommandContext) {
				enabled, _ := context.GetBooleanArg(0)
				if !context.IsError() {
					status := "disabled"
					if enabled {
						status = "enabled"
					}
					fmt.Printf("Hub DHCP authoritative successfully %s\n", status)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
