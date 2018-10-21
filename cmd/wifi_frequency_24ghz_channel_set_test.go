package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestWiFiFrequency24GhzChannelSetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().WiFiFrequency24GhzChannelSet(12345).Return(nil)

	AssertCommandOutput(t, NewWiFiFrequency24GhzChannelSetCommand(NewLoginCommand()), "12345")
}

func TestWiFiFrequency24GhzChannelSetCommandInvalidChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewWiFiFrequency24GhzChannelSetCommand(NewLoginCommand()), "abc")
}
