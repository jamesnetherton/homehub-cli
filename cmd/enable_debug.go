package cmd

import (
	"errors"
	"fmt"

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
				return
			}
			service.EnableDebug(enable)
		},
		PostExec: func(context *CommandContext) {
			if !context.IsError() {
				status := "disabled"
				enabled, _ := context.GetBooleanArg(0)
				if enabled {
					status = "enabled"
				}

				fmt.Printf("Hub client debugging %s\n", status)
			}
		},
	}
}
