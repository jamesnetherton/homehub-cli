package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
)

// AuthenticationRequiringCommand is a representation of a Hub CLI command that requires
// authentication before executing
type AuthenticationRequiringCommand struct {
	GenericCommand
	AuthenticatingCommand Command
}

// ExecuteLifecylce runs the command execution lifecycle
func (c *AuthenticationRequiringCommand) ExecuteLifecylce(args []string) {
	context := &CommandContext{
		args: args,
	}
	hub := service.GetHub()

	if (!helpRequested(context)) && (hub == nil || !service.IsLoggedIn()) {
		fmt.Printf("\nYou are not logged in. Please login...\n\n")
		c.AuthenticatingCommand.Execute(context)

		if context.GetResult() == nil {
			return
		}

		if context.IsError() || !context.GetResult().(bool) {
			fmt.Println("Login failed")
			return
		}
	}
	c.GenericCommand.ExecuteLifecylce(args)
}
