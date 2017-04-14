package cmd

import (
	"fmt"
	"testing"

	"github.com/jamesnetherton/homehub-cli/service"
)

type FakeCommand struct{}

func (c *FakeCommand) Execute(args []string) (result interface{}, err error) {
	return "Test result", nil
}

func (c *FakeCommand) ExecuteLifecylce(args []string) {
	c.Execute(args)
}

func (c *FakeCommand) GetName() string {
	return "TestCommand"
}

func (c *FakeCommand) Validate(args []string) bool {
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

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestCommandLineParseWithIncompleteFlag(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--username="}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestCommandLineParseWithIncompleteFlags(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--huburl=", "--username=", "--password="}
	commandLine := NewCommandLineParser(commands, args)

	expected := false
	actual, _ := commandLine.Parse()

	if expected != actual {
		t.Fatalf("Expected command line parse to return false")
	}
}

func TestCommandLineParseWithcompleteFlags(t *testing.T) {
	command := &FakeCommand{}
	commands := []Command{command}

	args := []string{"TestCommand", "--huburl=foo", "--username=bar", "--password=cheese"}
	commandLine := NewCommandLineParser(commands, args)

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
