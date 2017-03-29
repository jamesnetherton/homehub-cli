package homehub

import (
	"strconv"
	"strings"
)

var debug debugging

// Hub represents a Home Hub client for interacting with the router
type Hub struct {
	client *client
	URL    string
}

// New creates a new Hub client for interacting with the router
func New(URL string, username string, password string) *Hub {
	c := newClient(URL+"/cgi/json-req", username, hexmd5(password))
	return &Hub{c, URL}
}

// BandwidthMonitor returns bandwidth statistics for devices that have connected to the router
func (h *Hub) BandwidthMonitor() (result string, err error) {
	return h.client.getBandwidthUsage()
}

// BroadbandProductType returns the last used wan interface type. For BT this equates to the broadband product type
func (h *Hub) BroadbandProductType() (result string, err error) {
	i, err := h.client.getXPathValueString(mySagemcomBoxDeviceInfoInterfaceType)

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
func (h *Hub) ConnectedDevices() (result []DeviceDetail, err error) {
	return h.client.getXPathValues(mainCacheableHosts)
}

// DataPumpVersion returns the DSL line firmware version
func (h *Hub) DataPumpVersion() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoDatapumpVersion)
}

// DataReceived returns the the total data bytes received
func (h *Hub) DataReceived() (result int64, err error) {
	return h.client.getXPathValueInt64(mySagemcomBoxBasicStatusDataUsageReceived)
}

// DataSent returns the total data bytes sent
func (h *Hub) DataSent() (result int64, err error) {
	return h.client.getXPathValueInt64(mySagemcomBoxBasicStatusDataUsageSent)
}

// DeviceInfo returns infomation about a device matching the specified id
func (h *Hub) DeviceInfo(id int) (result DeviceDetail, err error) {
	host, err := h.client.getXPathHostValue(strings.Replace(ethernetDeviceDevicesList, "#", strconv.Itoa(id), 1))
	return host.DeviceDetail, err
}

// DhcpAuthoritative returns whether the hub is the authoritive DHCP server
func (h *Hub) DhcpAuthoritative() (result bool, err error) {
	return h.client.getXPathValueBool(mySagemcomBoxDhcpDhcpAuthoritative)
}

// DhcpPoolStart returns DHCP pool start IP adress
func (h *Hub) DhcpPoolStart() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDhcpIpv4PoolStart)
}

// DhcpPoolEnd returns DHCP pool end IP adress
func (h *Hub) DhcpPoolEnd() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDhcpIpv4PoolEnd)
}

// DhcpSubnetMask returns DHCP subnet mask
func (h *Hub) DhcpSubnetMask() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoLocalSubnetMask)
}

// DownstreamSyncSpeed returns the speed at which the router is downloading data
func (h *Hub) DownstreamSyncSpeed() (result int, err error) {
	return h.client.getXPathValueInt(mySagemcomBoxBasicStatusDownstreamSyncSpeedDsl)
}

// EnableDebug causes HTTP client request and responses to be output to the console
func (h *Hub) EnableDebug(enable bool) {
	if enable {
		debug = true
	} else {
		debug = false
	}
}

// EnableDhcpAuthoritative toggles whether the hub is the authoritive DHCP server
func (h *Hub) EnableDhcpAuthoritative(enable bool) (err error) {
	return h.client.setXPathValue(mySagemcomBoxDhcpDhcpAuthoritative, enable)
}

// EventLog returns the events that have taken place on the router since it was last reset
func (h *Hub) EventLog() (result *EventLog, err error) {
	e, err := h.client.getEventLog()

	if err != nil {
		return nil, err
	}

	eventLog := &EventLog{}

	for _, line := range strings.Split(e, "\n") {
		logEntry := parseEventLogEntry(line)
		if logEntry != nil {
			eventLog.Entries = append(eventLog.Entries, *logEntry)
		}
	}

	return eventLog, nil
}

// HardwareVersion returns the router hardware version
func (h *Hub) HardwareVersion() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoHardwareVersion)
}

// InternetConnectionStatus returns the internet connection status
func (h *Hub) InternetConnectionStatus() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoWanInternetStatus)
}

// LightBrightness returns the router LED brightness percentage
func (h *Hub) LightBrightness() (result int, err error) {
	return h.client.getXPathValueInt(mySagemcomBoxDeviceInfoHubLightBrightness)
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
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoHubLightStatus)
}

// LocalTime returns the local time from the router NTP server
func (h *Hub) LocalTime() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxMaintenanceNtpLocalTime)
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
	return h.client.getXPathValueString(technicalLogFirmwareVersion)
}

// PublicIPAddress returns the router public IP address
func (h *Hub) PublicIPAddress() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoPublicIpv4)
}

// PublicSubnetMask returns the router public subnet mask
func (h *Hub) PublicSubnetMask() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoPublicSubnetMask)
}

// Reboot restarts router serial number
func (h *Hub) Reboot() (err error) {
	// TODO: Figure out why the response is HTTP 500 and why there is a delay till the router reboots
	return h.client.doReboot()
}

// SambaIP returns the samba share IP address
func (h *Hub) SambaIP() (result string, err error) {
	return h.client.getXPathValueString(mymediaSambaIP)
}

// SambaHost returns a comma delimited list of samba host names
func (h *Hub) SambaHost() (result string, err error) {
	return h.client.getXPathValueString(mymediaSambaHost)
}

// SerialNumber returns the router serial number
func (h *Hub) SerialNumber() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoSerialNumber)
}

// SoftwareVersion returns the router software version
func (h *Hub) SoftwareVersion() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoSoftwareVersion)
}

// UpstreamSyncSpeed returns the speed at which the router is uploading data
func (h *Hub) UpstreamSyncSpeed() (result int, err error) {
	return h.client.getXPathValueInt(mySagemcomBoxBasicStatusUpstreamSyncSpeedDsl)
}

// Version returns the router version
func (h *Hub) Version() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoProductClass)
}

// WiFiSecurityMode returns the WiFi security mode in use
func (h *Hub) WiFiSecurityMode() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoWifi24SecurityMode)
}

// WiFiSSID returns the WiFi service set identifier
func (h *Hub) WiFiSSID() (result string, err error) {
	return h.client.getXPathValueString(mySagemcomBoxDeviceInfoWifi24Ssid)
}
