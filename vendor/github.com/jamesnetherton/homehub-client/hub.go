package homehub

import (
	"reflect"
	"strconv"
	"strings"
)

var debug debugging

// Hub represents a Home Hub client for interacting with the router
type Hub struct {
	client   *client
	URL      string
	firmware firmware
}

// New creates a new Hub client for interacting with the router
func New(URL string, username string, password string) *Hub {
	c := newClient(URL+"/cgi/json-req", username, password)
	return &Hub{c, URL, &firmwareSG4B1{}}
}

// BandwidthMonitor returns bandwidth statistics for devices that have connected to the router
func (h *Hub) BandwidthMonitor() (result *BandwidthLog, err error) {
	stats, err := h.client.getBandwidthUsage(h.firmware.bandwidthMonitorXPath())

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
	i, err := h.client.getXPathValueString(h.firmware.broadbandProductTypeXPath())

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
	devices, err := h.client.getXPathValues(h.firmware.connectedDevicesXPath(), reflect.TypeOf(d))

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
	return h.client.getXPathValueString(h.firmware.dataPumpVersionXPath())
}

// DataReceived returns the the total data bytes received
func (h *Hub) DataReceived() (result int64, err error) {
	return h.client.getXPathValueInt64(h.firmware.dataReceivedXPath())
}

// DataSent returns the total data bytes sent
func (h *Hub) DataSent() (result int64, err error) {
	return h.client.getXPathValueInt64(h.firmware.dataSentXPath())
}

// DeviceInfo returns infomation about a device matching the specified id
func (h *Hub) DeviceInfo(id int) (result *DeviceDetail, err error) {
	var hostType host
	valueType, err := h.client.getXPathValueType(strings.Replace(h.firmware.deviceInfoXPath(), "#", strconv.Itoa(id), 1), reflect.TypeOf(hostType))

	if err == nil {
		return &valueType.(*host).DeviceDetail, err
	}

	return nil, err
}

// DhcpAuthoritative returns whether the hub is the authoritive DHCP server
func (h *Hub) DhcpAuthoritative() (result bool, err error) {
	return h.client.getXPathValueBool(h.firmware.dhcpAuthoritativeXPath())
}

// DhcpPoolStart returns DHCP pool start IP adress
func (h *Hub) DhcpPoolStart() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.dhcpPoolStartXPath())
}

// DhcpPoolEnd returns DHCP pool end IP adress
func (h *Hub) DhcpPoolEnd() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.dhcpPoolEndXPath())
}

// DhcpSubnetMask returns DHCP subnet mask
func (h *Hub) DhcpSubnetMask() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.dhcpSubnetMaskXPath())
}

// DownstreamSyncSpeed returns the speed at which the router is downloading data
func (h *Hub) DownstreamSyncSpeed() (result int, err error) {
	return h.client.getXPathValueInt(h.firmware.downstreamSyncSpeedXPath())
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
	return h.client.setXPathValue(h.firmware.dhcpAuthoritativeXPath(), enable)
}

// EventLog returns the events that have taken place on the router since it was last reset
func (h *Hub) EventLog() (result *EventLog, err error) {
	e, err := h.client.getEventLog(h.firmware.eventLogXPath())

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
	return h.client.getXPathValueString(h.firmware.hardwareVersionXPath())
}

// InternetConnectionStatus returns the internet connection status
func (h *Hub) InternetConnectionStatus() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.internetConnectionStatusXPath())
}

// LightBrightness returns the router LED brightness percentage
func (h *Hub) LightBrightness() (result int, err error) {
	return h.client.getXPathValueInt(h.firmware.lightBrightnessXPath())
}

// LightBrightnessSet sets the router LED brightness percentage
func (h *Hub) LightBrightnessSet(brightness int) (err error) {
	return h.client.setXPathValue(h.firmware.lightBrightnessXPath(), brightness)
}

// LightEnable toggles the status of the router LED lights
func (h *Hub) LightEnable(enable bool) (err error) {
	status := "ON"
	if enable == false {
		status = "OFF"
	}
	return h.client.setXPathValue(h.firmware.lightEnableXPath(), status)
}

// LightStatus returns the router LED light satus
func (h *Hub) LightStatus() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.lightStatusXPath())
}

// LocalTime returns the local time from the router NTP server
func (h *Hub) LocalTime() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.localTimeXPath())
}

// Login authenticates a user
func (h *Hub) Login() (success bool, err error) {
	err = h.client.login()

	if err == nil {
		// Retrieve firmware version
		firmwareVersion, err := h.MaintenanceFirmwareVersion()
		if err == nil {
			if strings.HasPrefix(firmwareVersion, firmwareVersionSG4B1A) {
				h.firmware = &firmwareSG4B1A{}
			}
		} else {
			return false, err
		}

		return true, nil
	}

	return false, err
}

// MaintenanceFirmwareVersion returns the maintenance firmware version
func (h *Hub) MaintenanceFirmwareVersion() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.maintenanceFirmwareVersionXPath())
}

// NatRules returns IPV4 firewall NAT rules
func (h *Hub) NatRules() (result []NatRule, err error) {
	var r []NatRule
	natRules, err := h.client.getXPathValues(h.firmware.natRulesXPath(), reflect.TypeOf(r))

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
	valueType, err := h.client.getXPathValueType(strings.Replace(h.firmware.natRuleXPath(), "#", strconv.Itoa(id), 1), reflect.TypeOf(portMappingType))

	if err == nil {
		return &valueType.(*portMapping).NatRule, err
	}

	return nil, err
}

// NatRuleCreate creates an IPV4 firewall NAT rule
func (h *Hub) NatRuleCreate(natRule *NatRule) (err error) {
	uid, err := h.client.addChildXPathValue(h.firmware.natRuleCreateXPath(), &portMapping{NatRule: *natRule})

	if err == nil {
		natRule.UID = uid
		return nil
	}

	return err
}

// NatRuleDelete deletes an IPV4 firewall NAT rule
func (h *Hub) NatRuleDelete(id int) (err error) {
	return h.client.deleteChildXPathValue(strings.Replace(h.firmware.natRuleXPath(), "#", strconv.Itoa(id), 1))
}

// NatRuleUpdate updates an existing IPV4 firewall NAT rule
func (h *Hub) NatRuleUpdate(natRule NatRule) (err error) {
	return h.client.setXPathValues(natRule.getUpdateActions(h.firmware.natRuleXPath()))
}

// PublicIPAddress returns the router public IP address
func (h *Hub) PublicIPAddress() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.publicIPAddressXPath())
}

// PublicSubnetMask returns the router public subnet mask
func (h *Hub) PublicSubnetMask() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.publicSubnetMaskXPath())
}

// Reboot restarts the router
func (h *Hub) Reboot() (err error) {
	return h.client.doReboot(h.firmware.rebootXPath())
}

// SambaIP returns the samba share IP address
func (h *Hub) SambaIP() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.sambaIPXPath())
}

// SambaHost returns a comma delimited list of samba host names
func (h *Hub) SambaHost() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.sambaHostXPath())
}

// SerialNumber returns the router serial number
func (h *Hub) SerialNumber() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.serialNumberXPath())
}

// SoftwareVersion returns the router software version
func (h *Hub) SoftwareVersion() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.softwareVersionXPath())
}

// UpstreamSyncSpeed returns the speed at which the router is uploading data
func (h *Hub) UpstreamSyncSpeed() (result int, err error) {
	return h.client.getXPathValueInt(h.firmware.upstreamSyncSpeedXPath())
}

// Version returns the router version
func (h *Hub) Version() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.versionXPath())
}

// WiFiFrequency24Ghz returns information about the WiFI 2.4GHz frequency
func (h *Hub) WiFiFrequency24Ghz() (result *WiFiFrequency, err error) {
	var radioType radio
	valueType, err := h.client.getXPathValueType(h.firmware.wiFiFrequency24GhzXPath(), reflect.TypeOf(radioType))

	if err == nil {
		return &valueType.(*radio).WiFiFrequency, err
	}

	return nil, err
}

// WiFiFrequency24GhzChannelSet sets the operating channel for the 2.4Ghz frequency
func (h *Hub) WiFiFrequency24GhzChannelSet(channel int) (err error) {
	return h.client.setXPathValue(h.firmware.wiFiFrequency24GhzChannelSetXPath(), channel)
}

// WiFiFrequency5Ghz returns information about the WiFI 5GHz frequency
func (h *Hub) WiFiFrequency5Ghz() (result *WiFiFrequency, err error) {
	var radioType radio
	valueType, err := h.client.getXPathValueType(h.firmware.wiFiFrequency5GhzXPath(), reflect.TypeOf(radioType))

	if err == nil {
		return &valueType.(*radio).WiFiFrequency, err
	}

	return nil, err
}

// WiFiFrequency5GhzChannelSet sets the operating channel for the 5Ghz frequency
func (h *Hub) WiFiFrequency5GhzChannelSet(channel int) (err error) {
	return h.client.setXPathValue(h.firmware.wiFiFrequency5GhzChannelSetXPath(), channel)
}

// WiFiSecurityMode returns the WiFi security mode in use
func (h *Hub) WiFiSecurityMode() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.wiFiSecurityModeXPath())
}

// WiFiSSID returns the WiFi service set identifier
func (h *Hub) WiFiSSID() (result string, err error) {
	return h.client.getXPathValueString(h.firmware.wiFiSSIDXPath())
}
