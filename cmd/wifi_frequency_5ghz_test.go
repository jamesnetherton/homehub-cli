package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	"github.com/jamesnetherton/homehub-client"
)

func TestWiFiFrquency5GhzCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	frequency := &homehub.WiFiFrequency{
		UID:                12345,
		Alias:              "RADIO5G",
		Enable:             true,
		Status:             "UP",
		Channel:            54321,
		AvailableChannels:  "1,2,3,4",
		OperatingStandards: "a,n",
		SupportedStandards: "a,n",
	}

	hub.EXPECT().WiFiFrequency5Ghz().Return(frequency, nil)

	AssertCommandOutput(t, NewWiFiFrquency5GhzCommand(NewLoginCommand()))
}
