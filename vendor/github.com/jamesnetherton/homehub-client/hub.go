package homehub

import (
	"reflect"
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
func (h *Hub) BandwidthMonitor() (result *BandwidthLog, err error) {
	stats, err := h.client.getBandwidthUsage()

	if err != nil {
		return nil, err
	}

	bandwidthLog := &BandwidthLog{}

	for _, line := range strings.Split(stats, "\n") {
		logEntry := parseBandwidthLogEntry(line)
		if logEntry != nil {
			bandwidthLog.Entries = append(bandwidthLog.Entries, *logEntry)
		}
	}

	return bandwidthLog, nil
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
	var d []DeviceDetail
	devices, err := h.client.getXPathValues(mainCacheableHosts, reflect.TypeOf(d))

	if err == nil {
		var deviceDetails []DeviceDetail
		for _, deviceDetail := range devices {
			deviceDetails = append(deviceDetails, deviceDetail.(DeviceDetail))
		}
		return deviceDetails, nil
	}

	return nil, err
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
func (h *Hub) DeviceInfo(id int) (result *DeviceDetail, err error) {
	var hostType host
	valueType, err := h.client.getXPathValueType(strings.Replace(ethernetDeviceDevicesList, "#", strconv.Itoa(id), 1), reflect.TypeOf(hostType))

	if err == nil {
		return &valueType.(*host).DeviceDetail, err
	}

	return nil, err
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

// NatRules returns IPV4 firewall NAT rules
func (h *Hub) NatRules() (result []NatRule, err error) {
	var r []NatRule
	natRules, err := h.client.getXPathValues(accessControlPortForwardingPortmappings, reflect.TypeOf(r))

	if err == nil {
		var rules []NatRule
		for _, rule := range natRules {
			rules = append(rules, rule.(NatRule))
		}
		return rules, nil
	}

	return nil, err
}

// NatRule returns an IPV4 firewall NAT rule matching the specified id
func (h *Hub) NatRule(id int) (result *NatRule, err error) {
	var portMappingType portMapping
	valueType, err := h.client.getXPathValueType(strings.Replace(accessControlPortForwardingUID, "#", strconv.Itoa(id), 1), reflect.TypeOf(portMappingType))

	if err == nil {
		return &valueType.(*portMapping).NatRule, err
	}

	return nil, err
}

// NatRuleCreate creates an IPV4 firewall NAT rule
func (h *Hub) NatRuleCreate(natRule *NatRule) (err error) {
	uid, err := h.client.addChildXPathValue(accessControlPortForwardingPortmappings, &portMapping{NatRule: *natRule})

	if err == nil {
		natRule.UID = uid
		return nil
	}

	return err
}

// NatRuleDelete deletes an IPV4 firewall NAT rule
func (h *Hub) NatRuleDelete(id int) (err error) {
	return h.client.deleteChildXPathValue(strings.Replace(accessControlPortForwardingUID, "#", strconv.Itoa(id), 1))
}

// NatRuleUpdate updates an existing IPV4 firewall NAT rule
func (h *Hub) NatRuleUpdate(natRule NatRule) (err error) {
	return h.client.setXPathValues(natRule.getUpdateActions())
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
