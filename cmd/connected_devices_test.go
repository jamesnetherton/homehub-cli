package cmd

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jamesnetherton/homehub-cli/service"
	homehub "github.com/jamesnetherton/homehub-client"
)

func TestConnectedDevicesCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := NewMockHub(ctrl)

	var connectedDevices []homehub.DeviceDetail

	for i := 1; i <= 10; i++ {
		active := true
		interfaceType := "Ethernet"

		if i%2 == 0 {
			active = false
			interfaceType = "WiFi"
		} else if i%5 == 0 {
			interfaceType = "USB"
		}

		deviceDetail := &homehub.DeviceDetail{
			Active:          active,
			InterfaceType:   interfaceType,
			IPAddress:       fmt.Sprintf("192.168.1.%d", i),
			PhysicalAddress: fmt.Sprintf("AA:BB:CC:DD:EE:FF:0%d", i),
			UID:             i,
		}
		connectedDevices = append(connectedDevices, *deviceDetail)
	}

	service.SetHub(hub)
	service.AuthenticationComplete()

	hub.EXPECT().ConnectedDevices().Return(connectedDevices, nil)

	AssertCommandOutput(t, NewConnectedDevicesCommand(NewLoginCommand()))
}
