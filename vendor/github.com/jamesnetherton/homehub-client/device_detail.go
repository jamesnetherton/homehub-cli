package homehub

// DeviceDetail defines a device connected to the Home Hub
type DeviceDetail struct {
	UID                            int                `json:"uid,omitempty"`
	Alias                          string             `json:"Alias,omitempty"`
	PhysicalAddress                string             `json:"PhysAddress,omitempty"`
	IPAddress                      string             `json:"IPAddress,omitempty"`
	AddressSource                  string             `json:"AddressSource,omitempty"`
	DHCPClient                     string             `json:"DHCPClient,omitempty"`
	LeaseTimeRemaining             int                `json:"LeaseTimeRemaining,omitempty"`
	AssociatedDevice               string             `json:"AssociatedDevice,omitempty"`
	HostName                       string             `json:"HostName,omitempty"`
	Active                         bool               `json:"Active,omitempty"`
	LeaseStart                     int                `json:"LeaseStart,omitempty"`
	LeaseDuration                  int                `json:"LeaseDuration,omitempty"`
	InterfaceType                  string             `json:"InterfaceType,omitempty"`
	DetectedDeviceType             string             `json:"DetectedDeviceType,omitempty"`
	LastStateChange                string             `json:"LastStateChange,omitempty"`
	UserFriendlyName               string             `json:"UserFriendlyName,omitempty"`
	UserHostName                   string             `json:"UserHostName,omitempty"`
	UserDeviceType                 string             `json:"UserDeviceType,omitempty"`
	BlacklistEnable                bool               `json:"BlacklistEnable,omitempty"`
	UnblockHours                   int                `json:"UnblockHoursCount,omitempty"`
	Blacklisted                    bool               `json:"Blacklisted,omitempty"`
	BlacklistStatus                bool               `json:"BlacklistStatus,omitempty"`
	BlacklistedAccordingToSchedule bool               `json:"BlacklistedAccordingToSchedule,omitempty"`
	Hidden                         bool               `json:"Hidden,omitempty"`
	IPv4Addresses                  []IPAddress        `json:"IPv4Addresses,omitempty"`
	IPv6Addresses                  []IPAddress        `json:"IPv6Addresses,omitempty"`
	LastConnections                []ConnectionDetail `json:"LastConnections,omitempty"`
	ConnectionsAtLastReboot        int                `json:"ConnectionsNbreAtLastReboot,omitempty"`
}

type host struct {
	DeviceDetail `json:"Host,omitempty"`
}

// IPAddress defines an IPV4 or IPV6 address
type IPAddress struct {
	UID       int    `json:"uid,omitempty"`
	IPAddress string `json:"IPAddress,omitempty"`
	Active    bool   `json:"Active,omitempty"`
}

// ConnectionDetail defines connection information related to a device connected to the Home Hub
type ConnectionDetail struct {
	UID                 int    `json:"uid,omitempty"`
	ConnectionTimestamp string `json:"ConnectionTimestamp,omitempty"`
	DisconnectTimestamp string `json:"DisconnectionTimestamp,omitempty"`
}
