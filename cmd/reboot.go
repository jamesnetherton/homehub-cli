package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jamesnetherton/homehub-cli/service"
)

// NewRebootCommand creates a new command to invoke the Hub Reboot function
func NewRebootCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "Reboot",
			Description: "Reboots the Home Hub",
			Exec: func(context *CommandContext) {
				context.SetResult(nil, service.GetHub().Reboot())
			},
			PostExec: func(context *CommandContext) {
				fmt.Print("\nWaiting for Home Hub to reboot...")
				attempts := 0
				for {
					attempts++
					response, err := http.Get(service.GetHub().URL)
					if err != nil || response.StatusCode != 200 {
						if attempts == 24 {
							fmt.Println("\nGave up waiting for Home Hub to become available")
							break
						} else {
							fmt.Print(".")
						}
					} else {
						fmt.Println()
						break
					}
					time.Sleep(5000 * time.Millisecond)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
