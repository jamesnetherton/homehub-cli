package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestNatRuleCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	natRule := &homehub.NatRule{
		UID:                   12345,
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

	hub.EXPECT().NatRule(12345).Return(natRule, nil)

	AssertCommandOutput(t, NewNatRuleCommand(NewLoginCommand()), "12345")
}
