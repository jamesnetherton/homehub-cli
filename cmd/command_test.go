package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func AssertCommandOutput(t *testing.T, command Command, args ...string) {
	stdout := os.Stdout

	defer func() { os.Stdout = stdout }()

	fileName := fmt.Sprintf("%s.txt", t.Name())
	testOutput := filepath.Join(os.TempDir(), fileName)
	file, _ := os.Create(testOutput)
	os.Stdout = file

	command.ExecuteLifecylce(args)

	file.Close()

	expected, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Fatalf("Failed reading test data file %s", fileName)
	}

	actual, err := ioutil.ReadFile(testOutput)
	if err != nil {
		t.Fatalf("Failed reading test output file %s", testOutput)
	}

	if string(expected) != string(actual) {
		t.Errorf("\nExpected command output:\n\n%s\n\nActual command output:\n\n%s", expected, actual)
	}
}
