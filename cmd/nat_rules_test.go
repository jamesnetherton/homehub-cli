package cmd

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestNatRulesCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	var natRules []homehub.NatRule

	for i := 1; i <= 5; i++ {

		natRule := &homehub.NatRule{
			UID:                   i,
			Enable:                true,
			Alias:                 "",
			ExternalInterface:     "",
			AllExternalInterfaces: false,
			LeaseDuration:         0,
			RemoteHost:            "",
			ExternalPort:          i * 1000,
			ExternalPortEndRange:  i * 2000,
			InternalInterface:     "",
			InternalPort:          i * 3000,
			Protocol:              "TCP",
			Service:               "NONE",
			InternalClient:        fmt.Sprintf("192.168.1.%d", i),
			Description:           fmt.Sprintf("Test NAT rule %d", i),
			Creator:               "USER",
			Target:                "ACCEPT",
			LeaseStart:            "",
		}

		natRules = append(natRules, *natRule)
	}

	hub.EXPECT().NatRules().Return(natRules, nil)

	AssertCommandOutput(t, NewNatRulesCommand(NewLoginCommand()))
}
