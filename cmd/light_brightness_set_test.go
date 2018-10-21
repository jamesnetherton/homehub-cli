package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestLightBrightnessSetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().LightBrightnessSet(75).Return(nil)

	AssertCommandOutput(t, NewLightBrightnessSetCommand(NewLoginCommand()), "75")
}

func TestLightBrightnessSetCommandBelowMinimumBound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewLightBrightnessSetCommand(NewLoginCommand()), "-1")
}

func TestLightBrightnessSetCommandAboveMaximumBound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewLightBrightnessSetCommand(NewLoginCommand()), "101")
}

func TestLightBrightnessSetCommandInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewLightBrightnessSetCommand(NewLoginCommand()), "abc")
}
