package cmd

import (
	"testing"
)

func TestAboutCommand(t *testing.T) {
	AssertCommandOutput(t, NewAboutCommand())
}
