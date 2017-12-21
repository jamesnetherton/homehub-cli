package cmd

import (
	"errors"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewEnableDebugCommand creates a new command to invoke the Hub EnableDebug function
func NewEnableDebugCommand() *GenericCommand {
	return &GenericCommand{
		Name:        "EnableDebug",
		Description: "Enables debug logging of HTTP requests",
		ArgNames:    []string{"enable"},
		ArgTypes:    []string{"bool"},
		Exec: func(context *CommandContext) {
			enable, err := context.GetBooleanArg(0)
			if err != nil {
				parseErr := errors.New("Enable flag must be either true or false")
				context.SetResult(nil, parseErr)
			}
			service.GetHub().EnableDebug(enable)
		},
	}
}
