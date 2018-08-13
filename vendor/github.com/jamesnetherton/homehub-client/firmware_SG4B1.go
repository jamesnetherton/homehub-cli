package homehub

type firmwareSG4B1 struct {
}

func (f *firmwareSG4B1) bandwidthMonitorXPath() string {
	return bandwidthMonitoring
}

func (f *firmwareSG4B1) broadbandProductTypeXPath() string {
	return mySagemcomBoxDeviceInfoInterfaceType
}

func (f *firmwareSG4B1) connectedDevicesXPath() string {
	return mainCacheableHosts
}

func (f *firmwareSG4B1) dataPumpVersionXPath() string {
	return mySagemcomBoxDeviceInfoDatapumpVersion
}

func (f *firmwareSG4B1) dataReceivedXPath() string {
	return mySagemcomBoxBasicStatusDataUsageReceived
}

func (f *firmwareSG4B1) dataSentXPath() string {
	return mySagemcomBoxBasicStatusDataUsageSent
}

func (f *firmwareSG4B1) deviceInfoXPath() string {
	return ethernetDeviceDevicesList
}

func (f *firmwareSG4B1) dhcpAuthoritativeXPath() string {
	return mySagemcomBoxDhcpDhcpAuthoritative
}

func (f *firmwareSG4B1) dhcpPoolStartXPath() string {
	return mySagemcomBoxDhcpIpv4PoolStart
}

func (f *firmwareSG4B1) dhcpPoolEndXPath() string {
	return mySagemcomBoxDhcpIpv4PoolEnd
}

func (f *firmwareSG4B1) dhcpSubnetMaskXPath() string {
	return mySagemcomBoxDeviceInfoLocalSubnetMask
}

func (f *firmwareSG4B1) downstreamSyncSpeedXPath() string {
	return technicalLogDataRateDown
}

func (f *firmwareSG4B1) eventLogXPath() string {
	return eventLog
}

func (f *firmwareSG4B1) hardwareVersionXPath() string {
	return mySagemcomBoxDeviceInfoHardwareVersion
}

func (f *firmwareSG4B1) internetConnectionStatusXPath() string {
	return mySagemcomBoxDeviceInfoWanInternetStatus
}

func (f *firmwareSG4B1) lightBrightnessXPath() string {
	return mySagemcomBoxDeviceInfoHubLightBrightness
}

func (f *firmwareSG4B1) lightEnableXPath() string {
	return mySagemcomBoxDeviceInfoHubLightStatus
}

func (f *firmwareSG4B1) lightStatusXPath() string {
	return mySagemcomBoxDeviceInfoHubLightStatus
}

func (f *firmwareSG4B1) localTimeXPath() string {
	return mySagemcomBoxMaintenanceNtpLocalTime
}

func (f *firmwareSG4B1) maintenanceFirmwareVersionXPath() string {
	return technicalLogFirmwareVersion
}

func (f *firmwareSG4B1) natRulesXPath() string {
	return accessControlPortForwardingPortmappings
}

func (f *firmwareSG4B1) natRuleXPath() string {
	return accessControlPortForwardingUID
}

func (f *firmwareSG4B1) natRuleCreateXPath() string {
	return accessControlPortForwardingPortmappings
}

func (f *firmwareSG4B1) publicIPAddressXPath() string {
	return mySagemcomBoxDeviceInfoPublicIpv4
}

func (f *firmwareSG4B1) publicSubnetMaskXPath() string {
	return mySagemcomBoxDeviceInfoPublicSubnetMask
}

func (f *firmwareSG4B1) rebootXPath() string {
	return device
}

func (f *firmwareSG4B1) sambaIPXPath() string {
	return mymediaSambaIP
}

func (f *firmwareSG4B1) sambaHostXPath() string {
	return mymediaSambaHost
}

func (f *firmwareSG4B1) serialNumberXPath() string {
	return mySagemcomBoxDeviceInfoSerialNumber
}

func (f *firmwareSG4B1) softwareVersionXPath() string {
	return mySagemcomBoxDeviceInfoSoftwareVersion
}

func (f *firmwareSG4B1) upstreamSyncSpeedXPath() string {
	return technicalLogDataRateUp
}

func (f *firmwareSG4B1) versionXPath() string {
	return mySagemcomBoxDeviceInfoProductClass
}

func (f *firmwareSG4B1) wiFiFrequency24GhzXPath() string {
	return mySagemcomBoxDeviceInfoWifi24
}

func (f *firmwareSG4B1) wiFiFrequency24GhzChannelSetXPath() string {
	return technicalLogWifiChannel24
}

func (f *firmwareSG4B1) wiFiFrequency5GhzXPath() string {
	return mySagemcomBoxDeviceInfoWifi5O
}

func (f *firmwareSG4B1) wiFiFrequency5GhzChannelSetXPath() string {
	return technicalLogWifiChannel5
}

func (f *firmwareSG4B1) wiFiSecurityModeXPath() string {
	return mySagemcomBoxDeviceInfoWifi24SecurityMode
}

func (f *firmwareSG4B1) wiFiSSIDXPath() string {
	return mySagemcomBoxDeviceInfoWifi24Ssid
}
