package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const path = "../build/homehub-cli"

func readTestDataFile(fileName string) string {
	bytesRead, err := ioutil.ReadFile("testdata/" + fileName)
	if err == nil {
		return string(bytesRead)
	}
	return ""
}

func assertCliOutput(cliInput string, stubResponseFile string, t *testing.T) {
	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Stat(path + ".exe")
		if err != nil {
			t.Fatal("homehub-cli executable not present. Run 'make build'")
		}
	}

	response := readTestDataFile(stubResponseFile)
	if response == "" {
		t.Fatalf("Failed reading test data stubed response file %s", stubResponseFile)
	}

	header := readTestDataFile("banner.txt")
	expected := header + response

	cli := exec.Command(path)
	command := cliInput + "\n"
	cli.Stdin = strings.NewReader(command)

	outputBytes, _ := cli.Output()
	actual := string(outputBytes)

	if expected != actual {
		t.Errorf("\nExpected response:\n\n%s\n\nActual response:\n\n%s", expected, actual)
	}

	cli.Wait()
}

func TestCLIWithInvalidCommand(t *testing.T) {
	assertCliOutput("testing", "invalid_command.txt", t)
}

func TestCommandWithoutRequiredArguments(t *testing.T) {
	assertCliOutput("EnableDebug", "missing_required_args.txt", t)
}

func TestCommandHelpMessage(t *testing.T) {
	assertCliOutput("EnableDebug --help", "command_help.txt", t)
}

func TestCommandHelpMessageAuthenticatedCommand(t *testing.T) {
	assertCliOutput("BandwidthMonitor --help", "authenticated_command_help.txt", t)
}
