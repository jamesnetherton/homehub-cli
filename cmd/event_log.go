package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewEventLogCommand creates a new command to invoke the Hub EventLog function
func NewEventLogCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "EventLog",
			Description: "Gets the Home Hub event log entries",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().EventLog()) },
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					headerPattern := "%-30s%-20s%-25s%-7s\n"
					dataPattern := "%-30s%-20s%-25s%-7s\n"
					eventLog := context.GetResult().(*homehub.EventLog)
					eventLogEntries := eventLog.Entries

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")
					fmt.Printf(headerPattern, "Time", "Type", "Category", "Message")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")

					for i := 0; i < len(eventLogEntries); i++ {
						fmt.Printf(dataPattern,
							eventLogEntries[i].Timestamp,
							eventLogEntries[i].Type,
							eventLogEntries[i].Category,
							eventLogEntries[i].Message,
						)
					}
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
