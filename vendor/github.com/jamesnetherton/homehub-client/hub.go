package homehub

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Hub represents a Home Hub client for interacting with the router
type Hub struct {
	client *client
}

// New creates a new Hub client for interacting with the router
func New(URL string, username string, password string) *Hub {
	c := newClient(URL+"/cgi/json-req", username, password)
	log.SetPrefix("INFO: ")
	log.SetFlags(log.LstdFlags)
	log.SetOutput(ioutil.Discard)
	return &Hub{c}
}

// BroadbandProductType returns the last used wan interface type. For BT this equates to the broadband product type
func (h *Hub) BroadbandProductType() (result string, err error) {
	i, err := h.client.sendXPathRequest(mySagemcomBoxDeviceInfoInterfaceType)

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

// DataPumpVersion returns the DSL line firmware version
func (h *Hub) DataPumpVersion() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoDatapumpVersion)
}

// DataReceived returns the the total data bytes received
func (h *Hub) DataReceived() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxBasicStatusDataUsageReceived)
}

// DataSent returns the total data bytes sent
func (h *Hub) DataSent() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxBasicStatusDataUsageSent)
}

// DhcpPoolStart returns DHCP pool start IP adress
func (h *Hub) DhcpPoolStart() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDhcpIpv4PoolStart)
}

// DhcpPoolEnd returns DHCP pool end IP adress
func (h *Hub) DhcpPoolEnd() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDhcpIpv4PoolEnd)
}

// DhcpSubnetMask returns DHCP subnet mask
func (h *Hub) DhcpSubnetMask() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoLocalSubnetMask)
}

// DownstreamSyncSpeed returns the speed at which the router is downloading data
func (h *Hub) DownstreamSyncSpeed() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxBasicStatusDownstreamSyncSpeedDsl)
}

// EnableDebug causes HTTP client request and responses to be output to the console
func (h *Hub) EnableDebug(enable bool) {
	if enable {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

// HardwareVersion returns the router hardware version
func (h *Hub) HardwareVersion() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoHardwareVersion)
}

// InternetConnectionStatus returns the internet connection status
func (h *Hub) InternetConnectionStatus() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoWanInternetStatus)
}

// LightStatus returns the router LED light satus
func (h *Hub) LightStatus() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoHubLightStatus)
}

// LocalTime returns the local time from the router NTP server
func (h *Hub) LocalTime() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxMaintenanceNtpLocalTime)
}

// Login authenticates a user
func (h *Hub) Login() (success bool, err error) {
	req := newRequest(&h.client.authData, "logIn", "")
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
	return h.client.sendXPathRequest(technicalLogFirmwareVersion)
}

// PublicIPAddress returns the router public IP address
func (h *Hub) PublicIPAddress() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoPublicIpv4)
}

// PublicSubnetMask returns the router public subnet mask
func (h *Hub) PublicSubnetMask() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoPublicSubnetMask)
}

// Reboot restarts router serial number
func (h *Hub) Reboot() (result string, err error) {
	// TODO: Figure out why the response is HTTP 500 and why there is a delay till the router reboots
	return h.client.sendXPathRequest("Device")
}

// SambaIP returns the samba share IP address
func (h *Hub) SambaIP() (result string, err error) {
	return h.client.sendXPathRequest(mymediaSambaIP)
}

// SambaHost returns a comma delimited list of samba host names
func (h *Hub) SambaHost() (result string, err error) {
	return h.client.sendXPathRequest(mymediaSambaHost)
}

// SerialNumber returns the router serial number
func (h *Hub) SerialNumber() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoSerialNumber)
}

// SoftwareVersion returns the router software version
func (h *Hub) SoftwareVersion() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoSoftwareVersion)
}

// UpstreamSyncSpeed returns the speed at which the router is uploading data
func (h *Hub) UpstreamSyncSpeed() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxBasicStatusUpstreamSyncSpeedDsl)
}

// Version returns the router version
func (h *Hub) Version() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoProductClass)
}

// WiFiSecurityMode returns the WiFi security mode in use
func (h *Hub) WiFiSecurityMode() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoWifi24SecurityMode)
}

// WiFiSSID returns the WiFi service set identifier
func (h *Hub) WiFiSSID() (result string, err error) {
	return h.client.sendXPathRequest(mySagemcomBoxDeviceInfoWifi24Ssid)
}
