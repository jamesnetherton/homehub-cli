package homehub

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// Hub represents a Home Hub client for interacting with the router
type Hub struct {
	client *client
	URL    string
}

// New creates a new Hub client for interacting with the router
func New(URL string, username string, password string) *Hub {
	c := newClient(URL+"/cgi/json-req", username, hexmd5(password))
	log.SetPrefix("INFO: ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(ioutil.Discard)
	return &Hub{c, URL}
}

// BandwidthMonitor returns bandwidth statistics for devices that have connected to the router
func (h *Hub) BandwidthMonitor() (result string, err error) {
	bandwidthMonitorRequest := newBandwidthMonitorRequest(&h.client.authData)
	req := newHubResourceRequest(&h.client.authData, h.URL, bandwidthMonitorRequest)
	resp, err := req.send()
	if err != nil {
		return "", err
	}

	return resp.body, nil
}

// BroadbandProductType returns the last used wan interface type. For BT this equates to the broadband product type
func (h *Hub) BroadbandProductType() (result string, err error) {
	i, err := h.client.getXPathValue(mySagemcomBoxDeviceInfoInterfaceType)

	if err != nil {
		return "", err
	}

	if i == "FTTX" {
		return "BT Infinity 3 and 4", nil
	} else if i == "VDSL" {
		return "BT Infinity", nil
	}
	return "BT Broadband", nil
}

// ConnectedDevices returns information about any devices connected to the router
func (h *Hub) ConnectedDevices() (result string, err error) {
	values, err := h.client.getXPathValues(mainCacheableHosts)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	headerPattern := "%-5s%-20s%-25s%-7s\n"
	dataPattern := "%-5d%-20s%-25s%-7s\n"

	for _, value := range values {
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf(headerPattern, "--", "----------", "----------------", "----"))
		buffer.WriteString(fmt.Sprintf(headerPattern, "ID", "IP Address", "Physical Address", "Type"))
		buffer.WriteString(fmt.Sprintf(headerPattern, "--", "----------", "----------------", "----"))
		for _, device := range value {
			if device.InterfaceType == "WiFi" || device.InterfaceType == "Ethernet" {
				buffer.WriteString(fmt.Sprintf(dataPattern,
					device.UID,
					device.IPAddress,
					device.PhysicalAddress,
					device.InterfaceType,
				))
			}
		}
	}

	return buffer.String(), err
}

// DataPumpVersion returns the DSL line firmware version
func (h *Hub) DataPumpVersion() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoDatapumpVersion)
}

// DataReceived returns the the total data bytes received
func (h *Hub) DataReceived() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxBasicStatusDataUsageReceived)
}

// DataSent returns the total data bytes sent
func (h *Hub) DataSent() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxBasicStatusDataUsageSent)
}

// DeviceInfo returns infomation about a device matching the specified id
func (h *Hub) DeviceInfo(id int) (result string, err error) {
	host, err := h.client.getXPathHostValue(strings.Replace(ethernetDeviceDevicesList, "#", strconv.Itoa(id), 0))

	if err != nil {
		return "", nil
	}

	var buffer bytes.Buffer
	headerPattern := "%-5s%-20s%-25s%-7s\n"
	dataPattern := "%-5d%-20s%-25s%-7s\n"

	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf(headerPattern, "--", "----------", "----------------", "----"))
	buffer.WriteString(fmt.Sprintf(headerPattern, "ID", "IP Address", "Physical Address", "Type"))
	buffer.WriteString(fmt.Sprintf(headerPattern, "--", "----------", "----------------", "----"))
	buffer.WriteString(fmt.Sprintf(dataPattern,
		host.UID,
		host.IPAddress,
		host.PhysicalAddress,
		host.InterfaceType,
	))

	return buffer.String(), nil
}

// DhcpAuthoritative returns whether the hub is the authoritive DHCP server
func (h *Hub) DhcpAuthoritative() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDhcpDhcpAuthoritative)
}

// DhcpPoolStart returns DHCP pool start IP adress
func (h *Hub) DhcpPoolStart() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDhcpIpv4PoolStart)
}

// DhcpPoolEnd returns DHCP pool end IP adress
func (h *Hub) DhcpPoolEnd() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDhcpIpv4PoolEnd)
}

// DhcpSubnetMask returns DHCP subnet mask
func (h *Hub) DhcpSubnetMask() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoLocalSubnetMask)
}

// DownstreamSyncSpeed returns the speed at which the router is downloading data
func (h *Hub) DownstreamSyncSpeed() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxBasicStatusDownstreamSyncSpeedDsl)
}

// EnableDebug causes HTTP client request and responses to be output to the console
func (h *Hub) EnableDebug(enable bool) {
	if enable {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

// EnableDhcpAuthoritative toggles whether the hub is the authoritive DHCP server
func (h *Hub) EnableDhcpAuthoritative(enable bool) (err error) {
	return h.client.setXPathValue(mySagemcomBoxDhcpDhcpAuthoritative, enable)
}

// EventLog returns the events that have taken place on the router since it was last reset
func (h *Hub) EventLog() (result string, err error) {
	eventLogRequest := newEventLogRequest(&h.client.authData)
	req := newHubResourceRequest(&h.client.authData, h.URL, eventLogRequest)
	resp, err := req.send()
	if err != nil {
		return "", err
	}

	return resp.body, nil
}

// HardwareVersion returns the router hardware version
func (h *Hub) HardwareVersion() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoHardwareVersion)
}

// InternetConnectionStatus returns the internet connection status
func (h *Hub) InternetConnectionStatus() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoWanInternetStatus)
}

// LightBrightness returns the router LED brightness percentage
func (h *Hub) LightBrightness() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoHubLightBrightness)
}

// LightBrightnessSet sets the router LED brightness percentage
func (h *Hub) LightBrightnessSet(brightness int) (err error) {
	return h.client.setXPathValue(mySagemcomBoxDeviceInfoHubLightBrightness, brightness)
}

// LightEnable toggles the status of the router LED lights
func (h *Hub) LightEnable(enable bool) (err error) {
	status := "ON"
	if enable == false {
		status = "OFF"
	}
	return h.client.setXPathValue(mySagemcomBoxDeviceInfoHubLightStatus, status)
}

// LightStatus returns the router LED light satus
func (h *Hub) LightStatus() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoHubLightStatus)
}

// LocalTime returns the local time from the router NTP server
func (h *Hub) LocalTime() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxMaintenanceNtpLocalTime)
}

// Login authenticates a user
func (h *Hub) Login() (success bool, err error) {
	req := newLoginRequest(&h.client.authData)
	resp, err := req.send()

	if err == nil {
		h.client.authData.sessionID = strconv.Itoa(resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters.ID)
		h.client.authData.nonce = resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters.Nonce
		return true, nil
	}

	return false, err
}

// MaintenaceFirmwareVersion returns the maintenance firmware version
func (h *Hub) MaintenaceFirmwareVersion() (result string, err error) {
	return h.client.getXPathValue(technicalLogFirmwareVersion)
}

// PublicIPAddress returns the router public IP address
func (h *Hub) PublicIPAddress() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoPublicIpv4)
}

// PublicSubnetMask returns the router public subnet mask
func (h *Hub) PublicSubnetMask() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoPublicSubnetMask)
}

// Reboot restarts router serial number
func (h *Hub) Reboot() (result string, err error) {
	// TODO: Figure out why the response is HTTP 500 and why there is a delay till the router reboots
	return h.client.getXPathValue("Device")
}

// SambaIP returns the samba share IP address
func (h *Hub) SambaIP() (result string, err error) {
	return h.client.getXPathValue(mymediaSambaIP)
}

// SambaHost returns a comma delimited list of samba host names
func (h *Hub) SambaHost() (result string, err error) {
	return h.client.getXPathValue(mymediaSambaHost)
}

// SerialNumber returns the router serial number
func (h *Hub) SerialNumber() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoSerialNumber)
}

// SoftwareVersion returns the router software version
func (h *Hub) SoftwareVersion() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoSoftwareVersion)
}

// UpstreamSyncSpeed returns the speed at which the router is uploading data
func (h *Hub) UpstreamSyncSpeed() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxBasicStatusUpstreamSyncSpeedDsl)
}

// Version returns the router version
func (h *Hub) Version() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoProductClass)
}

// WiFiSecurityMode returns the WiFi security mode in use
func (h *Hub) WiFiSecurityMode() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoWifi24SecurityMode)
}

// WiFiSSID returns the WiFi service set identifier
func (h *Hub) WiFiSSID() (result string, err error) {
	return h.client.getXPathValue(mySagemcomBoxDeviceInfoWifi24Ssid)
}
