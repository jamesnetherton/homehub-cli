package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	"github.com/jamesnetherton/homehub-cli/util"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewConnectedDevicesCommand creates a new command to invoke the Hub ConnectedDevices function
func NewConnectedDevicesCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "ConnectedDevices",
			Description: "Gets details related to the devices connected to the Home Hub",
			Exec:        func(context *CommandContext) { context.SetResult(service.GetHub().ConnectedDevices()) },
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					connectedDevices := context.GetResult().([]homehub.DeviceDetail)

					data := []string{
						"ID | IP Address | Physical Address | Type | Active",
						"",
					}

					for i := 0; i < len(connectedDevices); i++ {
						if connectedDevices[i].InterfaceType == "WiFi" || connectedDevices[i].InterfaceType == "Ethernet" {
							line := fmt.Sprintf("%d | %s | %s | %s | %s", connectedDevices[i].UID, connectedDevices[i].IPAddress, connectedDevices[i].PhysicalAddress, connectedDevices[i].InterfaceType, util.HumanizeBool(connectedDevices[i].Active))
							data = append(data, line)
						}
					}

					util.Columnize(data)
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
