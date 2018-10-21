package cmd

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestDeviceInfoCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	deviceDetail := &homehub.DeviceDetail{
		Active:          true,
		InterfaceType:   "Ethernet",
		IPAddress:       "192.168.1.111",
		PhysicalAddress: "AA:BB:CC:DD:EE:FF:01",
		UID:             12345,
	}

	hub.EXPECT().DeviceInfo(12345).Return(deviceDetail, nil)

	AssertCommandOutput(t, NewDeviceInfoCommand(NewLoginCommand()), "12345")
}

func TestDeviceInfoCommandInvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)
	service.SetHub(hub)
	service.AuthenticationComplete()

	AssertCommandOutput(t, NewDeviceInfoCommand(NewLoginCommand()), "abc")
}
