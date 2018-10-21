package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestWiFiFrequency5GhzChannelSetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().WiFiFrequency5GhzChannelSet(12345).Return(nil)

	AssertCommandOutput(t, NewWiFiFrequency5GhzChannelSetCommand(NewLoginCommand()), "12345")
}

func TestWiFiFrequency5GhzChannelSetCommandInvalidChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewWiFiFrequency5GhzChannelSetCommand(NewLoginCommand()), "abc")
}
