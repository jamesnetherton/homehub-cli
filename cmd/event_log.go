package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	"github.com/jamesnetherton/homehub-cli/util"
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
					eventLog := context.GetResult().(*homehub.EventLog)
					eventLogEntries := eventLog.Entries

					data := []string{
						"Time | Type | Category | Message",
						"",
					}

					for i := 0; i < len(eventLogEntries); i++ {
						line := fmt.Sprintf("%s | %s | %s | %s", eventLogEntries[i].Timestamp, eventLogEntries[i].Type, eventLogEntries[i].Category, eventLogEntries[i].Message)
						data = append(data, line)
					}

					util.Columnize(data)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
