package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestLightEnableCommandEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().LightEnable(true).Return(nil)

	AssertCommandOutput(t, NewLightEnableCommand(NewLoginCommand()), "true")
}

func TestLightEnableCommandDisabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().LightEnable(false).Return(nil)

	AssertCommandOutput(t, NewLightEnableCommand(NewLoginCommand()), "false")
}

func TestLightEnableCommandInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewLightEnableCommand(NewLoginCommand()), "abc")
}
