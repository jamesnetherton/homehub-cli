package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestEnableDebugCommandEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().EnableDebug(true)

	AssertCommandOutput(t, NewEnableDebugCommand(), "true")
}

func TestEnableDebugCommandDisabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().EnableDebug(false)

	AssertCommandOutput(t, NewEnableDebugCommand(), "false")
}

func TestEnableDebugCommandInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewEnableDebugCommand(), "abc")
}

func TestEnableDebugWithoutLogin(t *testing.T) {
	AssertCommandOutput(t, NewEnableDebugCommand(), "true")
}
