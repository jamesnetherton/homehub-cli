package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
			fmt.Println(err)
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
		Exec: func(args []string) (result interface{}, err error) {
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
			return success, err
		},
	}

	enableDebug := &cmd.GenericCommand{
		Name:        "EnableDebug",
		Description: "Gets details about a specific device connected to the Home Hub",
		ArgNames:    []string{"enable"},
		ArgTypes:    []string{"bool"},
		Exec: func(args []string) (result interface{}, err error) {
			enable, err := strconv.ParseBool(args[0])
			if err != nil {
				parseErr := errors.New("Enable flag must be either true or false")
				return nil, parseErr
			}
			service.GetHub().EnableDebug(enable)
			return nil, nil
		},
	}
	bandwidthMonitor := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "BandwidthMonitor",
			Description: "Creates a new Home Hub login session",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().BandwidthMonitor() },
		},
		AuthenticatingCommand: login,
	}

	broadbandProductType := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "BroadbandProductType",
			Description: "Gets the BT broadband product type",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().BroadbandProductType() },
		},
		AuthenticatingCommand: login,
	}
	connectedDevices := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "ConnectedDevices",
			Description: "Gets details related to the devices connected to the Home Hub",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().ConnectedDevices() },
			PostExec: func(result interface{}, err error) error {
				if err == nil {
					headerPattern := "%-5s%-20s%-25s%-7s\n"
					dataPattern := "%-5d%-20s%-25s%-7s\n"
					connectedDevices := result.([]homehub.DeviceDetail)

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
					return nil
				}
				return err
			},
		},
		AuthenticatingCommand: login,
	}
	dataPumpVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataPumpVersion",
			Description: "Gets details related to the DSL line firmware version",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DataPumpVersion() },
		},
		AuthenticatingCommand: login,
	}
	dataReceived := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataReceived",
			Description: "Gets the number of bytes receieved since the Home Hub was last rebooted",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DataReceived() },
		},
		AuthenticatingCommand: login,
	}
	dataSent := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DataSent",
			Description: "Gets the number of bytes sent since the Home Hub was last rebooted",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DataSent() },
		},
		AuthenticatingCommand: login,
	}
	deviceInfo := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DeviceInfo",
			Description: "Gets details about a specific device connected to the Home Hub",
			ArgNames:    []string{"deviceId"},
			ArgTypes:    []string{"int"},
			Exec: func(args []string) (result interface{}, err error) {
				id, err := strconv.Atoi(args[0])
				if err != nil {
					parseErr := errors.New("Device ID must be a numeric value")
					return nil, parseErr
				}
				return service.GetHub().DeviceInfo(id)
			},
			PostExec: func(result interface{}, err error) error {
				if err == nil {
					headerPattern := "%-5s%-20s%-25s%-7s\n"
					dataPattern := "%-5d%-20s%-25s%-7s\n"
					device := result.(homehub.DeviceDetail)

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
					return nil
				}
				return err
			}},
		AuthenticatingCommand: login,
	}
	dhcpAuthoritative := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpAuthoritative",
			Description: "Gets details about whether the Home Hub is the authoritative DHCP server",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DhcpAuthoritative() },
		},
		AuthenticatingCommand: login,
	}
	dhcpPoolEnd := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpPoolEnd",
			Description: "Gets the Home Hub IPV4 DHCP pool end address",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DhcpPoolEnd() },
		},
		AuthenticatingCommand: login,
	}
	dhcpPoolStart := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpPoolStart",
			Description: "Gets the Home Hub IPV4 DHCP pool start address",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DhcpPoolStart() },
		},
		AuthenticatingCommand: login,
	}
	dhcpSubnetMask := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DhcpSubnetMask",
			Description: "Gets the Home Hub DHCP subnet mask",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DhcpSubnetMask() },
		},
		AuthenticatingCommand: login,
	}
	downstreamSyncSpeed := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "DownstreamSyncSpeed",
			Description: "Gets the available speed at which the Home Hub can download data",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().DownstreamSyncSpeed() },
		},
		AuthenticatingCommand: login,
	}
	enableDhcpAuthoritative := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "EnableDhcpAuthoritative",
			Description: "Toggles whether the Home Hub is the authoritative DHCP server",
			ArgNames:    []string{"enable"},
			ArgTypes:    []string{"bool"},
			Exec: func(args []string) (result interface{}, err error) {
				enable, err := strconv.ParseBool(args[0])
				if err != nil {
					parseErr := errors.New("Enable flag must be either true or false")
					return nil, parseErr
				}
				err = service.GetHub().EnableDhcpAuthoritative(enable)
				return nil, err
			},
		},
		AuthenticatingCommand: login,
	}
	eventLog := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "EventLog",
			Description: "Gets the Home Hub event log entries",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().EventLog() },
			PostExec: func(result interface{}, err error) error {
				if err == nil {
					headerPattern := "%-30s%-20s%-25s%-7s\n"
					dataPattern := "%-30s%-20s%-25s%-7s\n"
					eventLog := result.(*homehub.EventLog)
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
					return nil
				}
				return err
			},
		},
		AuthenticatingCommand: login,
	}
	hardwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "HardwareVersion",
			Description: "Gets the Home Hub hardware version",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().HardwareVersion() },
		},
		AuthenticatingCommand: login,
	}
	internetConnectionStatus := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "InternetConnectionStatus",
			Description: "Gets the status of the Home Hub internet connection",
			Exec: func(args []string) (result interface{}, err error) {
				return service.GetHub().InternetConnectionStatus()
			},
		},
		AuthenticatingCommand: login,
	}
	lightBrightness := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightBrightness",
			Description: "Gets the Home Hub LED brightness percentage value",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().LightBrightness() },
		},
		AuthenticatingCommand: login,
	}
	lightBrightnessSet := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightBrightnessSet",
			Description: "Sets the Home Hub LED brightness percentage value",
			Exec: func(args []string) (result interface{}, err error) {
				brightness, _ := strconv.Atoi(args[0])
				err = service.GetHub().LightBrightnessSet(brightness)
				return nil, err
			},
		},
		AuthenticatingCommand: login,
	}
	lightEnable := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightEnable",
			Description: "Toggles the Home Hub LED on or off",
			Exec: func(args []string) (result interface{}, err error) {
				enable, _ := strconv.ParseBool(args[0])
				err = service.GetHub().LightEnable(enable)
				return nil, err
			},
		},
		AuthenticatingCommand: login,
	}
	lightStatus := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LightStatus",
			Description: "Gets the status of the Home Hub LED",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().LightStatus() },
		},
		AuthenticatingCommand: login,
	}
	localTime := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "LocalTime",
			Description: "Gets local time from the Home Hub",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().LocalTime() },
		},
		AuthenticatingCommand: login,
	}
	maintenaceFirmwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "MaintenaceFirmwareVersion",
			Description: "Gets the Home Hub maintenance firmware version",
			Exec: func(args []string) (result interface{}, err error) {
				return service.GetHub().MaintenaceFirmwareVersion()
			},
		},
		AuthenticatingCommand: login,
	}
	publicIPAddress := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "PublicIPAddress",
			Description: "Gets the Home Hub public IP address",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().PublicIPAddress() },
		},
		AuthenticatingCommand: login,
	}
	publicSubnetMask := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "PublicSubnetMask",
			Description: "Gets the Home Hub public IP subnet mask",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().PublicSubnetMask() },
		},
		AuthenticatingCommand: login,
	}
	reboot := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "Reboot",
			Description: "Reboots the Home Hub",
			Exec: func(args []string) (result interface{}, err error) {
				return nil, service.GetHub().Reboot()
			},
			PostExec: func(result interface{}, err error) error {
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
				return nil
			},
		},
		AuthenticatingCommand: login,
	}
	sambaHost := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SambaHost",
			Description: "Gets the Home Hub samba host name",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().SambaHost() },
		},
		AuthenticatingCommand: login,
	}
	sambaIP := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SambaIP",
			Description: "Gets the Home Hub samba IP address",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().SambaIP() },
		},
		AuthenticatingCommand: login,
	}
	serialNumber := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SerialNumber",
			Description: "Gets the Home Hub serial number",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().SerialNumber() },
		},
		AuthenticatingCommand: login,
	}
	softwareVersion := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "SoftwareVersion",
			Description: "Gets the Home Hub software version",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().SoftwareVersion() },
		},
		AuthenticatingCommand: login,
	}
	upstreamSyncSpeed := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "UpstreamSyncSpeed",
			Description: "Gets the Home Hub upload speed",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().UpstreamSyncSpeed() },
		},
		AuthenticatingCommand: login,
	}
	version := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "Version",
			Description: "Gets the Home Hub version",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().Version() },
		},
		AuthenticatingCommand: login,
	}
	wifiSSID := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "WiFiSSID",
			Description: "Gets the Home Hub WiFI SSID",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().WiFiSSID() },
		},
		AuthenticatingCommand: login,
	}
	wifiSecurityMode := &cmd.AuthenticationRequiringCommand{
		GenericCommand: cmd.GenericCommand{
			Name:        "WiFiSecurityMode",
			Description: "Gets the Home Hub WiFI security mode",
			Exec:        func(args []string) (result interface{}, err error) { return service.GetHub().WiFiSecurityMode() },
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
