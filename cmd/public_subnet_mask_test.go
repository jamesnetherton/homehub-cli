package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
)

func TestPublicSubnetMaskCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().PublicSubnetMask().Return("11.0.0.0", nil)

	AssertCommandOutput(t, NewPublicSubnetMaskCommand(NewLoginCommand()))
}
