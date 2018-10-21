package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestDhcpSubnetMaskCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().DhcpSubnetMask().Return("192.0.0.0", nil)

	AssertCommandOutput(t, NewDhcpSubnetMaskCommand(NewLoginCommand()))
}
