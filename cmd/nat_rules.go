package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	"github.com/jamesnetherton/homehub-cli/util"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewNatRulesCommand creates a new command to invoke the Hub NatRules function
func NewNatRulesCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "NatRules",
			Description: "Gets any IPV4 NAT rules configured on the Home Hub",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().NatRules()) },
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					natRules := context.GetResult().([]homehub.NatRule)

					data := []string{
						"ID | Description | Enabled | External Port Start | External Port End | Internal Port Start | Internal Port End | Protocol",
						"",
					}

					for i := 0; i < len(natRules); i++ {
						line := fmt.Sprintf("%d | %s | %t | %d | %d | %d | %d | %s", natRules[i].UID, natRules[i].Description, natRules[i].Enable, natRules[i].ExternalPort, natRules[i].ExternalPortEndRange, natRules[i].InternalPort, natRules[i].ExternalPortEndRange, natRules[i].Protocol)
						data = append(data, line)
					}

					util.Columnize(data)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
