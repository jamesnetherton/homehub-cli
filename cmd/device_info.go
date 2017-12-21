package cmd

import (
	"errors"
	"fmt"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewDeviceInfoCommand creates a new command to invoke the Hub DeviceInfo function
func NewDeviceInfoCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "DeviceInfo",
			Description: "Gets details about a specific device connected to the Home Hub",
			ArgNames:    []string{"deviceId"},
			ArgTypes:    []string{"int"},
			Exec: func(context *CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("Device ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(service.GetHub().DeviceInfo(id))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					headerPattern := "%-5s%-20s%-25s%-7s\n"
					dataPattern := "%-5d%-20s%-25s%-7s\n"
					device := context.GetResult().(*homehub.DeviceDetail)

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")
					fmt.Printf(headerPattern, "ID", "IP Address", "Physical Address", "Type")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")

					fmt.Printf(dataPattern,
						device.UID,
						device.IPAddress,
						device.PhysicalAddress,
						device.InterfaceType,
					)
				}
			}},
		AuthenticatingCommand: authenticatingCommand,
	}
}
