package cmd

import (
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
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
					headerPattern := "%-5s%-20s%-25s%-20s%-5s\n"
					dataPattern := "%-5d%-20s%-25s%-20s%-5s\n"
					connectedDevices := context.GetResult().([]homehub.DeviceDetail)

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----", "------")
					fmt.Printf(headerPattern, "ID", "IP Address", "Physical Address", "Type", "Active")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----", "------")

					for i := 0; i < len(connectedDevices); i++ {
						if connectedDevices[i].InterfaceType == "WiFi" || connectedDevices[i].InterfaceType == "Ethernet" {
							fmt.Printf(dataPattern,
								connectedDevices[i].UID,
								connectedDevices[i].IPAddress,
								connectedDevices[i].PhysicalAddress,
								connectedDevices[i].InterfaceType,
								humanizeBool(connectedDevices[i].Active),
							)
						}
					}
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}
