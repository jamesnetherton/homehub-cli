package cmd

import (
	"fmt"
	"testing"

	"github.com/jamesnetherton/homehub-cli/service"
)

type FakeCommand struct {
}

func (c *FakeCommand) Execute(context *CommandContext) {
	context.SetResult("Test result", nil)
}

func (c *FakeCommand) ExecuteLifecylce(args []string) {
	context := &CommandContext{
		args: args,
	}

	c.Execute(context)
}

func (c *FakeCommand) GetName() string {
	return "TestCommand"
}

func (c *FakeCommand) Validate(context *CommandContext) bool {
	return true
}

func (c *FakeCommand) Usage() {
	fmt.Printf("Usage FakeCommand foo<Bar>")
}

func TestCommandLineParseWithNoArgs(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestCommandLineParseWithOneArg(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand"}
	commandLine := NewCommandLineParser(commands, args)

	expected := true
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return true")
	}
}

func TestAuthenticatingCommandLineParseWithOneArg(t *testing.T) {
	command := &AuthenticationRequiringCommand{}

	// command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand"}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestAuthenticatingCommandLineParseWithIncompleteFlag(t *testing.T) {
	command := &AuthenticationRequiringCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--username="}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestAuthenticatingCommandLineParseWithIncompleteFlags(t *testing.T) {
	command := &AuthenticationRequiringCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--huburl=", "--username=", "--password="}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestAuthenticatingCommandLineParseWithcompleteFlags(t *testing.T) {
	command := &AuthenticationRequiringCommand{
		GenericCommand: GenericCommand{
			Name: "TestCommand",
			Exec: func(context *CommandContext) {
			},
		},
		AuthenticatingCommand: &GenericCommand{
			Exec: func(context *CommandContext) {
				context.SetResult(true, nil)
			},
			PostExec: func(context *CommandContext) {},
		},
	}
	commands := []Command{command}

	args := []string{"TestCommand", "--huburl=foo", "--username=bar", "--password=cheese"}
	commandLine := NewCommandLineParser(commands, args)

	service.NewHub("http://fake.url.for.testing", "fake username", "fake password")
	service.AuthenticationComplete()

	expected := true
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return true")
	}
}

func TestCommandLineParseWithDefaultFlags(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--password=cheese"}
	commandLine := NewCommandLineParser(commands, args)

	service.AuthenticationComplete()

	expected := true
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return true")
	}
}
