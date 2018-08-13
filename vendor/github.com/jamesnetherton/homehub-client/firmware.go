package homehub

const (
	firmwareVersionSG4B1  string = "SG4B1"
	firmwareVersionSG4B1A string = firmwareVersionSG4B1 + "A"
)

type firmware interface {
	bandwidthMonitorXPath() string
	broadbandProductTypeXPath() string
	connectedDevicesXPath() string
	dataPumpVersionXPath() string
	dataReceivedXPath() string
	dataSentXPath() string
	deviceInfoXPath() string
	dhcpAuthoritativeXPath() string
	dhcpPoolStartXPath() string
	dhcpPoolEndXPath() string
	dhcpSubnetMaskXPath() string
	downstreamSyncSpeedXPath() string
	eventLogXPath() string
	hardwareVersionXPath() string
	internetConnectionStatusXPath() string
	lightBrightnessXPath() string
	lightEnableXPath() string
	lightStatusXPath() string
	localTimeXPath() string
	maintenanceFirmwareVersionXPath() string
	natRulesXPath() string
	natRuleXPath() string
	natRuleCreateXPath() string
	publicIPAddressXPath() string
	publicSubnetMaskXPath() string
	rebootXPath() string
	sambaIPXPath() string
	sambaHostXPath() string
	serialNumberXPath() string
	softwareVersionXPath() string
	upstreamSyncSpeedXPath() string
	versionXPath() string
	wiFiFrequency24GhzXPath() string
	wiFiFrequency24GhzChannelSetXPath() string
	wiFiFrequency5GhzXPath() string
	wiFiFrequency5GhzChannelSetXPath() string
	wiFiSecurityModeXPath() string
	wiFiSSIDXPath() string
}
