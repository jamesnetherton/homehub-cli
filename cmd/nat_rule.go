package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewNatRuleCommand creates a new command to invoke the Hub NatRule function
func NewNatRuleCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "NatRule",
			Description: "Gets an IPV4 NAT rule configured for the specified ID",
			ArgNames:    []string{"id"},
			ArgTypes:    []string{"int"},
			Exec: func(context *CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(service.GetHub().NatRule(id))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					dataPattern := "%-25s%-5d\n%-25s%-25s\n%-25s%-25s\n%-25s%-25s\n%-25s%-15t\n%-25s%-25s\n%-25s%-25d\n%-25s%-25d\n%-25s%-25s\n%-25s%-25d\n%-25s%-25d\n%-25s%-25s\n%-25s%-10s\n"
					natRule := context.GetResult().(*homehub.NatRule)

					fmt.Print("\n")
					fmt.Printf(dataPattern,
						"ID",
						natRule.UID,
						"Description",
						natRule.Description,
						"Alias",
						natRule.Alias,
						"Creator",
						natRule.Creator,
						"Enabled",
						natRule.Enable,
						"Service",
						natRule.Service,
						"External Port Start",
						natRule.ExternalPort,
						"External Port End",
						natRule.ExternalPortEndRange,
						"External IP",
						natRule.RemoteHost,
						"Internal Port Start",
						natRule.InternalPort,
						"Internal Port End",
						natRule.ExternalPortEndRange,
						"Internal IP",
						natRule.InternalClient,
						"Protocol",
						natRule.Protocol,
					)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
