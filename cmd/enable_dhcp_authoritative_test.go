package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestEnableDhcpAuthoritativeCommandEnabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().EnableDhcpAuthoritative(true)

	AssertCommandOutput(t, NewEnableDhcpAuthoritativeCommand(NewLoginCommand()), "true")
}

func TestEnableDhcpAuthoritativeCommandDisabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().EnableDhcpAuthoritative(false)

	AssertCommandOutput(t, NewEnableDhcpAuthoritativeCommand(NewLoginCommand()), "false")
}

func TestEnableDhcpAuthoritativeCommandInvalidValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewEnableDhcpAuthoritativeCommand(NewLoginCommand()), "abc")
}
