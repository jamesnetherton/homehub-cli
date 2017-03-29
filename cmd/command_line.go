package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesnetherton/homehub-cli/service"
)

// CommandLineParser is a representation of a command line parser engine
type CommandLineParser struct {
	Commands []Command
	Args     []string
}

// NewCommandLineParser creates a new CommandLineParser
func NewCommandLineParser(commands []Command, args []string) *CommandLineParser {
	return &CommandLineParser{
		commands,
		args,
	}
}

// Parse parses command line arguments
func (c *CommandLineParser) Parse() (result bool, err error) {
	if len(c.Args) == 0 {
		return false, errors.New("No command was specified")
	}

	commandName := c.Args[0]
	command := c.findMatchingCommand(commandName)

	if command != nil {
		fullCommandLine := strings.Join(c.Args, " ")

		re := regexp.MustCompile("(--.*)")
		args := re.FindAllString(fullCommandLine, -1)

		if len(args) == 0 {
			return false, errors.New("Home Hub login credentials are missing")
		}

		for _, match := range re.FindAllString(fullCommandLine, -1) {
			fullCommandLine = strings.Replace(fullCommandLine, match, "", -1)
		}

		hubURL, urlErr := c.getMandatoryArgument("huburl", args[0])
		if urlErr != nil {
			return false, errors.New("--huburl flag is missing")
		}

		userName, userNameErr := c.getMandatoryArgument("username", args[0])
		if userNameErr != nil {
			return false, errors.New("--username flag is missing")
		}

		password, passwordErr := c.getMandatoryArgument("password", args[0])
		if passwordErr != nil {
			return false, errors.New("--password flag is missing")
		}

		if service.IsLoggedIn() == false {
			hub := service.NewHub(hubURL, userName, password)
			success, err := hub.Login()

			if !success || err != nil {
				return false, errors.New("Login failed")
			}
		}

		service.AuthenticationComplete()
		fullCommandLine = strings.TrimSpace(strings.Replace(fullCommandLine, commandName, "", -1))

		if len(fullCommandLine) == 0 {
			command.ExecuteLifecylce([]string{})
		} else {
			command.ExecuteLifecylce(strings.Split(fullCommandLine, " "))
		}

		return true, nil
	}
	return false, errors.New("Unknown command: " + commandName)
}

func (c *CommandLineParser) getMandatoryArgument(argumentName string, argLine string) (argument string, err error) {
	args := strings.Split(argLine, " ")

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--"+argumentName) {
			argParts := strings.Split(args[i], "=")
			if len(argParts) != 2 {
				return "", errors.New("Expected " + argumentName + " to be in the format --" + argumentName + "=value")
			}
			return argParts[1], nil
		}
	}
	return "", errors.New("Did not find an argument named " + argumentName)
}

func (c *CommandLineParser) findMatchingCommand(commandName string) Command {
	for i := 0; i < len(c.Commands); i++ {
		if c.Commands[i].GetName() == commandName {
			return c.Commands[i]
		}
	}
	return nil
}

func (c *CommandLineParser) getCommandNames() []string {
	var commandNames []string
	for i := 0; i < len(c.Commands); i++ {
		if c.Commands[i].GetName() != "Login" && c.Commands[i].GetName() != "EnableDebug" {
			commandNames = append(commandNames, c.Commands[i].GetName())
		}
	}
	return commandNames
}

// PrintUsage prints the CLI usage message
func (c *CommandLineParser) PrintUsage() {
	fmt.Println("Usage:\n  homehub-cli [command] [args] --huburl=<home hub url> --username=<home hub username> --password=<home hub password>")
	fmt.Println("\nCommands:\n ", strings.Join(c.getCommandNames(), "\n  "))
}
