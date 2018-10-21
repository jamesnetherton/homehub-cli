package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestDhcpPoolStartCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().DhcpPoolStart().Return("192.168.1.0", nil)

	AssertCommandOutput(t, NewDhcpPoolStartCommand(NewLoginCommand()))
}
