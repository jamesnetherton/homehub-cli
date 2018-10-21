package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestWiFiSecurityModeCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().WiFiSecurityMode().Return("Test security mode", nil)

	AssertCommandOutput(t, NewWiFiSecurityModeCommand(NewLoginCommand()))
}
