package homehub

// WiFiFrequency represents a Home Hub wireless frequency
type WiFiFrequency struct {
	UID                int    `json:"uid"`
	Enable             bool   `json:"Enable"`
	Alias              string `json:"Alias"`
	Status             string `json:"Status"`
	SupportedStandards string `json:"SupportedStandards"`
	OperatingStandards string `json:"OperatingStandards"`
	AvailableChannels  string `json:"ChannelsInUse"`
	Channel            int    `json:"Channel"`
}

type radio struct {
	WiFiFrequency `json:"Radio,omitempty"`
}
