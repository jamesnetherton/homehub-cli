package cmd

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

// NewNatRuleCreateCommand creates a new command to invoke the Hub NatRuleCreate function
func NewNatRuleCreateCommand(authenticatingCommand *GenericCommand) *AuthenticationRequiringCommand {
	return &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name:        "NatRuleCreate",
			Description: "Creates an IPV4 NAT rule",
			ArgNames:    []string{"name", "ip-address", "external-port-start", "external-port-end", "internal-port-start", "protocol", "action"},
			ArgTypes:    []string{"string", "string", "int", "int", "int", "string", "string"},
			Exec: func(context *CommandContext) {
				name := context.GetStringArg(0)

				ip := net.ParseIP(context.GetStringArg(1))
				if ip == nil {
					parseErr := errors.New("ip address must be a valid IPV4 address")
					context.SetResult(nil, parseErr)
					return
				}

				externalPortStart, err := context.GetIntArg(2)
				if err != nil || externalPortStart <= 0 {
					parseErr := errors.New("External port start must be a positive numeric value")
					context.SetResult(nil, parseErr)
					return
				}

				externalPortEnd, err := context.GetIntArg(3)
				if err != nil || externalPortEnd <= 0 {
					parseErr := errors.New("External port end must be a positive numeric value")
					context.SetResult(nil, parseErr)
					return
				}

				internalPortStart, err := context.GetIntArg(4)
				if err != nil || internalPortStart <= 0 {
					parseErr := errors.New("Internal port start must be a positive numeric value")
					context.SetResult(nil, parseErr)
					return
				}

				protocol := strings.ToUpper(context.GetStringArg(5))
				if !isValidProtocol(protocol) {
					parseErr := errors.New("Protocol not supported. Valid values are TCP, UDP, BOTH")
					context.SetResult(nil, parseErr)
					return
				}

				action := strings.ToUpper(context.GetStringArg(6))
				if !isValidAction(action) {
					parseErr := errors.New("Action not supported. Valid values are DROP, ACCEPT, REJECT")
					context.SetResult(nil, parseErr)
					return
				}

				natRule := &homehub.NatRule{
					Enable:                true,
					Alias:                 "",
					ExternalInterface:     "",
					AllExternalInterfaces: false,
					LeaseDuration:         0,
					RemoteHost:            "",
					ExternalPort:          externalPortStart,
					ExternalPortEndRange:  externalPortEnd,
					InternalInterface:     "",
					InternalPort:          internalPortStart,
					Protocol:              protocol,
					Service:               "NONE",
					InternalClient:        ip.String(),
					Description:           name,
					Creator:               "USER",
					Target:                action,
					LeaseStart:            "",
				}

				context.SetResult(nil, service.GetHub().NatRuleCreate(natRule))
			},
			PostExec: func(context *CommandContext) {
				if !context.IsError() {
					fmt.Printf("NAT rule successfully created\n")
				}
			},
		},
		AuthenticatingCommand: authenticatingCommand,
	}
}

func isValidProtocol(protocol string) bool {
	protocols := [...]string{"TCP", "UDP", "BOTH"}
	for _, validProtocol := range protocols {
		if validProtocol == protocol {
			return true
		}
	}

	return false
}

func isValidAction(action string) bool {
	actions := [...]string{"DROP", "ACCEPT", "REJECT"}
	for _, validAction := range actions {
		if validAction == action {
			return true
		}
	}

	return false
}

