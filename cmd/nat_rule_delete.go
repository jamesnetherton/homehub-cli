package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewNatRuleDeleteCommand creates a new command to invoke the Hub NatRuleDelete function
func NewNatRuleDeleteCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "NatRuleDelete",
			Description: "Deletes an IPV4 NAT rule configured for the specified ID",
			ArgNames:    []string{"id"},
			ArgTypes:    []string{"int"},
			Exec: func(context *CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(nil, service.GetHub().NatRuleDelete(id))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					fmt.Printf("NAT rule successfully deleted\n")
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
