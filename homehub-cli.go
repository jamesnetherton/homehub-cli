package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/jamesnetherton/homehub-cli/cli"
	"github.com/jamesnetherton/homehub-cli/cmd"
	"github.com/jamesnetherton/homehub-cli/util"
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

	if !util.StringIsEmpty(os.Getenv("USER")) {
		user = os.Getenv("USER")
	} else if !util.StringIsEmpty(os.Getenv("USERNAME")) {
		user = os.Getenv("USERNAME")
	} else {
		user = "unknown"
	}

	return fmt.Sprintf("%s@homehub: ", user)
}

func initCommands() []cmd.Command {
	var commands []cmd.Command

	login := cmd.NewLoginCommand()

	commands = append(commands,
		login,
		cmd.NewAboutCommand(),
		cmd.NewBandwidthMonitorCommand(login),
		cmd.NewBroadbandProductTypeCommand(login),
		cmd.NewConnectedDevicesCommand(login),
		cmd.NewDataPumpVersionCommand(login),
		cmd.NewDataReceivedCommand(login),
		cmd.NewDataSentCommand(login),
		cmd.NewDeviceInfoCommand(login),
		cmd.NewDhcpAuthoritativeCommand(login),
		cmd.NewDhcpPoolEndCommand(login),
		cmd.NewDhcpPoolStartCommand(login),
		cmd.NewDhcpSubnetMaskCommand(login),
		cmd.NewDownstreamSyncSpeedCommand(login),
		cmd.NewEnableDebugCommand(),
		cmd.NewEnableDhcpAuthoritativeCommand(login),
		cmd.NewEventLogCommand(login),
		cmd.NewHardwareVersionCommand(login),
		cmd.NewInternetConnectionStatusCommand(login),
		cmd.NewLightBrightnessCommand(login),
		cmd.NewLightBrightnessSetCommand(login),
		cmd.NewLightEnableCommand(login),
		cmd.NewLightStatusCommand(login),
		cmd.NewLocalTimeCommand(login),
		cmd.NewMaintenanceFirmwareVersionCommand(login),
		cmd.NewNatRulesCommand(login),
		cmd.NewNatRuleCommand(login),
		cmd.NewNatRuleCreateCommand(login),
		cmd.NewNatRuleDeleteCommand(login),
		cmd.NewPublicIPAddressCommand(login),
		cmd.NewPublicSubnetMaskCommand(login),
		cmd.NewRebootCommand(login),
		cmd.NewSambaHostCommand(login),
		cmd.NewSambaIPCommand(login),
		cmd.NewSerialNumberCommand(login),
		cmd.NewSoftwareVersionCommand(login),
		cmd.NewUpstreamSyncSpeedCommand(login),
		cmd.NewVersionCommand(login),
		cmd.NewWiFiFrquency24GhzCommand(login),
		cmd.NewWiFiFrequency24GhzChannelSetCommand(login),
		cmd.NewWiFiFrquency5GhzCommand(login),
		cmd.NewWiFiFrequency5GhzChannelSetCommand(login),
		cmd.NewWiFiSSIDCommand(login),
		cmd.NewWiFiSecurityModeCommand(login))
	return commands
}
