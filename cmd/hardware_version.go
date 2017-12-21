package cmd

import "github.com/jamesnetherton/homehub-cli/service"

// NewHardwareVersionCommand creates a new command to invoke the Hub HardwareVersion function
func NewHardwareVersionCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "HardwareVersion",
			Description: "Gets the Home Hub hardware version",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().HardwareVersion()) },
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
