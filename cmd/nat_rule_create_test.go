package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestNatRuleCreateCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	natRule := &homehub.NatRule{
		Enable:                true,
		Alias:                 "",
		ExternalInterface:     "",
		AllExternalInterfaces: false,
		LeaseDuration:         0,
		RemoteHost:            "",
		ExternalPort:          1111,
		ExternalPortEndRange:  2222,
		InternalInterface:     "",
		InternalPort:          3333,
		Protocol:              "TCP",
		Service:               "NONE",
		InternalClient:        "192.168.1.100",
		Description:           "Test NAT rule",
		Creator:               "USER",
		Target:                "ACCEPT",
		LeaseStart:            "",
	}

	hub.EXPECT().NatRuleCreate(natRule).Return(nil)

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "1111", "2222", "3333", "TCP", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.ABC", "1111", "2222", "3333", "TCP", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidExternalPortStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "ABC", "2222", "3333", "TCP", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidExternalPortEnd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "1111", "ABC", "3333", "TCP", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidInternalPortStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "1111", "2222", "ABC", "TCP", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidProtocol(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "1111", "2222", "3333", "INVALID", "ACCEPT")
}

func TestNatRuleCreateCommandInvalidAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewNatRuleCreateCommand(NewLoginCommand()), "Test NAT rule", "192.168.1.100", "1111", "2222", "3333", "TCP", "INVALID")
}
