package cli

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const path = "../build/homehub-cli"

func TestCLIWithInvalidCommand(t *testing.T) {

	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	cli := exec.Command(path)
	command := "testing\n"
	cli.Stdin = strings.NewReader(command)

	output, _ := cli.Output()

	if !strings.Contains(string(output), "homehub: Command not found: \"testing\"") {
		t.Errorf("Expected command not found message for command '%s'", command)
	}

	cli.Wait()
}

func TestLoginPrompt(t *testing.T) {

	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	cli := exec.Command(path)
	command := "BandwidthMonitor\n"
	cli.Stdin = strings.NewReader(command)

	output, _ := cli.Output()

	if !strings.Contains(string(output), "You are not logged in. Please login") {
		t.Errorf("Expected please login message")
	}

	cli.Wait()
}

func TestCommandWithoutRequiredArguments(t *testing.T) {

	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	cli := exec.Command(path)
	command := "EnableDebug\n"
	cli.Stdin = strings.NewReader(command)

	output, _ := cli.Output()
	if !strings.Contains(string(output), "Usage: EnableDebug enable<bool>") {
		t.Errorf("Expected usage message: EnableDebug enable<bool>")
	}

	cli.Wait()
}

func TestCommandHelpMessage(t *testing.T) {

	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	cli := exec.Command(path)
	command := "EnableDebug --help\n"
	cli.Stdin = strings.NewReader(command)

	output, _ := cli.Output()
	if !strings.Contains(string(output), "EnableDebug: Enables debug logging of HTTP requests") {
		t.Errorf("Expected function help message")
	}

	cli.Wait()
}

func TestCommandHelpMessageAuthenticatedCommand(t *testing.T) {

	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	cli := exec.Command(path)
	command := "BandwidthMonitor --help\n"
	cli.Stdin = strings.NewReader(command)

	output, _ := cli.Output()
	if !strings.Contains(string(output), "BandwidthMonitor: Displays bandwidth statistics for devices that have connected to the Home Hub") {
		t.Errorf("Expected function help message")
	}

	cli.Wait()
}
