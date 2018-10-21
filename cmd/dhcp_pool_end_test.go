package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestDhcpPoolEndCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().DhcpPoolEnd().Return("192.168.1.99", nil)

	AssertCommandOutput(t, NewDhcpPoolEndCommand(NewLoginCommand()))
}
