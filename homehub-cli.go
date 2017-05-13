package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/chzyer/readline"
	"github.com/jamesnetherton/homehub-cli/cli"
	"github.com/jamesnetherton/homehub-cli/cmd"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func main() {

	commands := initCommands()

	if len(os.Args[1:]) == 0 {
		readLine, err := readline.NewEx(&readline.Config{
			Prompt:          initPrompt(),
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})

		if err != nil {
			panic(err)
		}

		cli := cli.NewCLI(commands, readLine)
		cli.Run()
	} else {
		commandLine := cmd.NewCommandLineParser(commands, os.Args[1:])
		success, err := commandLine.Parse()
		if !success {
			if err != nil {
				fmt.Printf("%s\n\n", err.Error())
			}
			commandLine.PrintUsage()
			os.Exit(1)
		}
	}
}

func initPrompt() string {
	var user string

	if !service.StringIsEmpty(os.Getenv("USER")) {
		user = os.Getenv("USER")
	} else if !service.StringIsEmpty(os.Getenv("USERNAME")) {
		user = os.Getenv("USERNAME")
	} else {
		user = "unknown"
	}

	return fmt.Sprintf("%s@homehub: ", user)
}

func initCommands() []cmd.Command {
	var commands []cmd.Command

	login := &cmd.GenericCommand{
		Name:        "Login",
		Description: "Creates a new Home Hub login session",
		Exec: func(context *cmd.CommandContext) {
			hub := service.GetHub()
			if hub == nil || !service.IsLoggedIn() {
				var hubURL string
				var userName string

				fmt.Print("Home hub URL: ")
				fmt.Scan(&hubURL)

				fmt.Print("Home hub user: ")
				fmt.Scan(&userName)

				password, _ := speakeasy.Ask(fmt.Sprint("Home Hub password: "))

				service.NewHub(hubURL, userName, password)
			}
			success, err := service.GetHub().Login()
			if success {
				service.AuthenticationComplete()
			}

			context.SetResult(success, err)
		},
	}

	enableDebug := &cmd.GenericCommand{
		Name:        "EnableDebug",
		Description: "Enables debug logging of HTTP requests",
		ArgNames:    []string{"enable"},
		ArgTypes:    []string{"bool"},
		Exec: func(context *cmd.CommandContext) {
			enable, err := context.GetBooleanArg(0)
			if err != nil {
				parseErr := errors.New("Enable flag must be either true or false")
				context.SetResult(nil, parseErr)
			}
			service.GetHub().EnableDebug(enable)
		},
	}
	bandwidthMonitor := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "BandwidthMonitor",
			Description: "Displays bandwidth statistics for devices that have connected to the Home Hub",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().BandwidthMonitor()) },
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					headerPattern := "%-30s%-20s%-25s%-7s\n"
					dataPattern := "%-30s%-20s%-25d%-7d\n"
					bandwidthLog := context.GetResult().(*homehub.BandwidthLog)
					bandwidthLogEntries := bandwidthLog.Entries

					fmt.Print("\n")
					fmt.Printf(headerPattern, "------------------", "----------", "----------", "--------")
					fmt.Printf(headerPattern, "   MAC Address    ", "   Date   ", "Downloaded", "Uploaded")
					fmt.Printf(headerPattern, "------------------", "----------", "----------", "--------")

					for i := 0; i < len(bandwidthLogEntries); i++ {
						fmt.Printf(dataPattern,
							bandwidthLogEntries[i].MACAddress,
							bandwidthLogEntries[i].Date,
							bandwidthLogEntries[i].DownloadMegabytes,
							bandwidthLogEntries[i].UploadMegabytes,
						)
					}

					fmt.Print("\n")
				}
			},
		},
		AuthenticatingCommand: login,
	}

	broadbandProductType := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "BroadbandProductType",
			Description: "Gets the BT broadband product type",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().BroadbandProductType()) },
		},
		AuthenticatingCommand: login,
	}
	connectedDevices := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "ConnectedDevices",
			Description: "Gets details related to the devices connected to the Home Hub",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().ConnectedDevices()) },
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					headerPattern := "%-5s%-20s%-25s%-7s\n"
					dataPattern := "%-5d%-20s%-25s%-7s\n"
					connectedDevices := context.GetResult().([]homehub.DeviceDetail)

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")
					fmt.Printf(headerPattern, "ID", "IP Address", "Physical Address", "Type")
					fmt.Printf(headerPattern, "--", "----------", "----------------", "----")

					for i := 0; i < len(connectedDevices); i++ {
						if connectedDevices[i].InterfaceType == "WiFI" || connectedDevices[i].InterfaceType == "Ethernet" {
							fmt.Printf(dataPattern,
								connectedDevices[i].UID,
								connectedDevices[i].IPAddress,
								connectedDevices[i].PhysicalAddress,
								connectedDevices[i].InterfaceType,
							)
						}
					}
				}
			},
		},
		AuthenticatingCommand: login,
	}
	dataPumpVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataPumpVersion",
			Description: "Gets details related to the DSL line firmware version",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DataPumpVersion()) },
		},
		AuthenticatingCommand: login,
	}
	dataReceived := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataReceived",
			Description: "Gets the number of bytes receieved since the Home Hub was last rebooted",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DataReceived()) },
		},
		AuthenticatingCommand: login,
	}
	dataSent := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataSent",
			Description: "Gets the number of bytes sent since the Home Hub was last rebooted",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DataSent()) },
		},
		AuthenticatingCommand: login,
	}
	deviceInfo := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DeviceInfo",
			Description: "Gets details about a specific device connected to the Home Hub",
			ArgNames:    []string{"deviceId"},
			ArgTypes:    []string{"int"},
			Exec: func(context *cmd.CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("Device ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(service.GetHub().DeviceInfo(id))
			},
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					headerPattern := "%-5s%-20s%-25s%-7s\n"
					dataPattern := "%-5d%-20s%-25s%-7s\n"
					device := context.GetResult().(homehub.DeviceDetail)

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
		AuthenticatingCommand: login,
	}
	dhcpAuthoritative := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpAuthoritative",
			Description: "Gets details about whether the Home Hub is the authoritative DHCP server",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DhcpAuthoritative()) },
		},
		AuthenticatingCommand: login,
	}
	dhcpPoolEnd := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpPoolEnd",
			Description: "Gets the Home Hub IPV4 DHCP pool end address",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DhcpPoolEnd()) },
		},
		AuthenticatingCommand: login,
	}
	dhcpPoolStart := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpPoolStart",
			Description: "Gets the Home Hub IPV4 DHCP pool start address",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DhcpPoolStart()) },
		},
		AuthenticatingCommand: login,
	}
	dhcpSubnetMask := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpSubnetMask",
			Description: "Gets the Home Hub DHCP subnet mask",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DhcpSubnetMask()) },
		},
		AuthenticatingCommand: login,
	}
	downstreamSyncSpeed := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DownstreamSyncSpeed",
			Description: "Gets the available speed at which the Home Hub can download data",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().DownstreamSyncSpeed()) },
		},
		AuthenticatingCommand: login,
	}
	enableDhcpAuthoritative := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "EnableDhcpAuthoritative",
			Description: "Toggles whether the Home Hub is the authoritative DHCP server",
			ArgNames:    []string{"enable"},
			ArgTypes:    []string{"bool"},
			Exec: func(context *cmd.CommandContext) {
				enable, err := context.GetBooleanArg(0)
				if err != nil {
					parseErr := errors.New("Enable flag must be either true or false")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(nil, service.GetHub().EnableDhcpAuthoritative(enable))
			},
		},
		AuthenticatingCommand: login,
	}
	eventLog := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "EventLog",
			Description: "Gets the Home Hub event log entries",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().EventLog()) },
			PostExec: func(context *cmd.CommandContext) {
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
		AuthenticatingCommand: login,
	}
	hardwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "HardwareVersion",
			Description: "Gets the Home Hub hardware version",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().HardwareVersion()) },
		},
		AuthenticatingCommand: login,
	}
	internetConnectionStatus := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "InternetConnectionStatus",
			Description: "Gets the status of the Home Hub internet connection",
			Exec: func(context *cmd.CommandContext) {
				context.SetResult(service.GetHub().InternetConnectionStatus())
			},
		},
		AuthenticatingCommand: login,
	}
	lightBrightness := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightBrightness",
			Description: "Gets the Home Hub LED brightness percentage value",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().LightBrightness()) },
		},
		AuthenticatingCommand: login,
	}
	lightBrightnessSet := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightBrightnessSet",
			Description: "Sets the Home Hub LED brightness percentage value",
			Exec: func(context *cmd.CommandContext) {
				brightness, _ := context.GetIntArg(0)
				context.SetResult(nil, service.GetHub().LightBrightnessSet(brightness))
			},
		},
		AuthenticatingCommand: login,
	}
	lightEnable := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightEnable",
			Description: "Toggles the Home Hub LED on or off",
			Exec: func(context *cmd.CommandContext) {
				enable, _ := context.GetBooleanArg(0)
				context.SetResult(nil, service.GetHub().LightEnable(enable))
			},
		},
		AuthenticatingCommand: login,
	}
	lightStatus := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightStatus",
			Description: "Gets the status of the Home Hub LED",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().LightStatus()) },
		},
		AuthenticatingCommand: login,
	}
	localTime := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LocalTime",
			Description: "Gets local time from the Home Hub",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().LocalTime()) },
		},
		AuthenticatingCommand: login,
	}
	maintenaceFirmwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "MaintenaceFirmwareVersion",
			Description: "Gets the Home Hub maintenance firmware version",
			Exec: func(context *cmd.CommandContext) {
				context.SetResult(service.GetHub().MaintenaceFirmwareVersion())
			},
		},
		AuthenticatingCommand: login,
	}
	natRules := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "NatRules",
			Description: "Gets any IPV4 NAT rules configured on the Home Hub",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().NatRules()) },
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					headerPattern := "%-5s%-25s%-15s%-25s%-25s%-25s%-25s%-10s\n"
					dataPattern := "%-5d%-25s%-15t%-25d%-25d%-25d%-25d%-10s\n"
					natRules := context.GetResult().([]homehub.NatRule)

					fmt.Print("\n")
					fmt.Printf(headerPattern, "--", "---------------------", "-------", "-------------------", "-----------------", "-------------------", "-----------------", "--------")
					fmt.Printf(headerPattern, "ID", "Description", "Enabled", "External Port Start", "External Port End", "Internal Port Start", "Internal Port End", "Protocol")
					fmt.Printf(headerPattern, "--", "---------------------", "-------", "-------------------", "-----------------", "-------------------", "-----------------", "--------")

					for i := 0; i < len(natRules); i++ {
						fmt.Printf(dataPattern,
							natRules[i].UID,
							natRules[i].Description,
							natRules[i].Enable,
							natRules[i].ExternalPort,
							natRules[i].ExternalPortEndRange,
							natRules[i].InternalPort,
							natRules[i].ExternalPortEndRange,
							natRules[i].Protocol,
						)
					}
				}
			},
		},
		AuthenticatingCommand: login,
	}
	natRule := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "NatRule",
			Description: "Gets an IPV4 NAT rule configured for the specified ID",
			ArgNames:    []string{"id"},
			ArgTypes:    []string{"int"},
			Exec: func(context *cmd.CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(service.GetHub().NatRule(id))
			},
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					dataPattern := "%-25s%-5d\n%-25s%-25s\n%-25s%-25s\n%-25s%-25s\n%-25s%-15t\n%-25s%-25s\n%-25s%-25d\n%-25s%-25d\n%-25s%-25s\n%-25s%-25d\n%-25s%-25d\n%-25s%-25s\n%-25s%-10s\n"
					natRule := context.GetResult().(*homehub.NatRule)

					fmt.Print("\n")
					fmt.Printf(dataPattern,
						"ID",
						natRule.UID,
						"Description",
						natRule.Description,
						"Alias",
						natRule.Alias,
						"Creator",
						natRule.Creator,
						"Enabled",
						natRule.Enable,
						"Service",
						natRule.Service,
						"External Port Start",
						natRule.ExternalPort,
						"External Port End",
						natRule.ExternalPortEndRange,
						"External IP",
						natRule.RemoteHost,
						"Internal Port Start",
						natRule.InternalPort,
						"Internal Port End",
						natRule.ExternalPortEndRange,
						"Internal IP",
						natRule.InternalClient,
						"Protocol",
						natRule.Protocol,
					)
				}
			},
		},
		AuthenticatingCommand: login,
	}
	natRuleDelete := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "NatRuleDelete",
			Description: "Deletes an IPV4 NAT rule configured for the specified ID",
			ArgNames:    []string{"id"},
			ArgTypes:    []string{"int"},
			Exec: func(context *cmd.CommandContext) {
				id, err := context.GetIntArg(0)
				if err != nil {
					parseErr := errors.New("ID must be a numeric value")
					context.SetResult(nil, parseErr)
				}
				context.SetResult(nil, service.GetHub().NatRuleDelete(id))
			},
			PostExec: func(context *cmd.CommandContext) {
				if !context.IsError() {
					fmt.Printf("NAT rule successfully deleted\n")
				}
			},
		},
		AuthenticatingCommand: login,
	}
	publicIPAddress := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "PublicIPAddress",
			Description: "Gets the Home Hub public IP address",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().PublicIPAddress()) },
		},
		AuthenticatingCommand: login,
	}
	publicSubnetMask := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "PublicSubnetMask",
			Description: "Gets the Home Hub public IP subnet mask",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().PublicSubnetMask()) },
		},
		AuthenticatingCommand: login,
	}
	reboot := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "Reboot",
			Description: "Reboots the Home Hub",
			Exec: func(context *cmd.CommandContext) {
				context.SetResult(nil, service.GetHub().Reboot())
			},
			PostExec: func(context *cmd.CommandContext) {
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
		AuthenticatingCommand: login,
	}
	sambaHost := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SambaHost",
			Description: "Gets the Home Hub samba host name",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().SambaHost()) },
		},
		AuthenticatingCommand: login,
	}
	sambaIP := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SambaIP",
			Description: "Gets the Home Hub samba IP address",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().SambaIP()) },
		},
		AuthenticatingCommand: login,
	}
	serialNumber := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SerialNumber",
			Description: "Gets the Home Hub serial number",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().SerialNumber()) },
		},
		AuthenticatingCommand: login,
	}
	softwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SoftwareVersion",
			Description: "Gets the Home Hub software version",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().SoftwareVersion()) },
		},
		AuthenticatingCommand: login,
	}
	upstreamSyncSpeed := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "UpstreamSyncSpeed",
			Description: "Gets the Home Hub upload speed",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().UpstreamSyncSpeed()) },
		},
		AuthenticatingCommand: login,
	}
	version := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "Version",
			Description: "Gets the Home Hub version",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().Version()) },
		},
		AuthenticatingCommand: login,
	}
	wifiSSID := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "WiFiSSID",
			Description: "Gets the Home Hub WiFI SSID",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().WiFiSSID()) },
		},
		AuthenticatingCommand: login,
	}
	wifiSecurityMode := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "WiFiSecurityMode",
			Description: "Gets the Home Hub WiFI security mode",
			Exec:        func(context *cmd.CommandContext) { context.SetResult(service.GetHub().WiFiSecurityMode()) },
		},
		AuthenticatingCommand: login,
	}

	commands = append(commands,
		login,
		enableDebug,
		bandwidthMonitor,
		broadbandProductType,
		connectedDevices,
		dataPumpVersion,
		dataReceived,
		dataSent,
		deviceInfo,
		dhcpAuthoritative,
		dhcpPoolEnd,
		dhcpPoolStart,
		dhcpSubnetMask,
		downstreamSyncSpeed,
		enableDhcpAuthoritative,
		eventLog,
		hardwareVersion,
		internetConnectionStatus,
		lightBrightness,
		lightBrightnessSet,
		lightEnable,
		lightStatus,
		localTime,
		maintenaceFirmwareVersion,
		natRules,
		natRule,
		natRuleDelete,
		publicIPAddress,
		publicSubnetMask,
		reboot,
		sambaHost,
		sambaIP,
		serialNumber,
		softwareVersion,
		upstreamSyncSpeed,
		version,
		wifiSSID,
		wifiSecurityMode)
	return commands
}
