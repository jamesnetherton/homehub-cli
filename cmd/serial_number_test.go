package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestSerialNumberCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().SerialNumber().Return("1234567890", nil)

	AssertCommandOutput(t, NewSerialNumberCommand(NewLoginCommand()))
}
