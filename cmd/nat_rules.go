package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
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
					headerPattern := "%-5s%-25s%-15s%-25s%-25s%-25s%-25s%-10s\n"
					dataPattern := "%-5d%-25s%-15t%-25d%-25d%-25d%-25d%-10s\n"
					natRules := context.GetResult().([]homehub.NatRule)

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "---------------------", "-------", "-------------------", "-----------------", "-------------------", "-----------------", "--------")
					fmt.Printf(headerPattern, "ID", "Description", "Enabled", "External Port Start", "External Port End", "Internal Port Start", "Internal Port End", "Protocol")
					fmt.Printf(headerPattern, "--", "---------------------", "-------", "-------------------", "-----------------", "-------------------", "-----------------", "--------")

					for i := 0; i < len(natRules); i++ {
						fmt.Printf(dataPattern,
							natRules[i].UID,
							natRules[i].Description,
							natRules[i].Enable,
							natRules[i].ExternalPort,
							natRules[i].ExternalPortEndRange,
							natRules[i].InternalPort,
							natRules[i].ExternalPortEndRange,
							natRules[i].Protocol,
						)
					}
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
