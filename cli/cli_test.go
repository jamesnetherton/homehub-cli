package cli

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIWithInvalidCommand(t *testing.T) {

	path := "../build/homehub-cli"

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

	path := "../build/homehub-cli"

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

	path := "../build/homehub-cli"

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
