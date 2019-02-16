package service

import (
	"time"

	homehub "github.com/jamesnetherton/homehub-client"
)

var hub homehub.Hub
var url string
var isLoggedIn bool
var ticker *time.Ticker
var debug bool

// NewHub creates a new Hub
func NewHub(hubURL string, userName string, password string) homehub.Hub {
	url = hubURL
	hub = homehub.New(hubURL, userName, password)
	hub.EnableDebug(debug)

	// Poll the hub every minute to see if the user session is still active
	if ticker == nil {
		ticker := time.NewTicker(time.Millisecond * 60000)
		go func() {
			for range ticker.C {
				_, err := hub.Version()
				if err == nil {
					isLoggedIn = true
				} else {
					isLoggedIn = false
				}
			}
		}()
	}

	return hub
}

// GetHub gets the stored instance of Hub
func GetHub() homehub.Hub {
	return hub
}

// SetHub sets the stored instance of Hub
func SetHub(hubInstance homehub.Hub) {
	hub = hubInstance
}

// GetHubURL gets the Hub URL
func GetHubURL() string {
	return url
}

// IsLoggedIn returns whether a user is logged into the Hub
func IsLoggedIn() bool {
	return isLoggedIn
}

// AuthenticationComplete flags that the Hub login process is complete
func AuthenticationComplete() {
	isLoggedIn = true
}

// EnableDebug enables or disabled debugging
func EnableDebug(enabled bool) {
	debug = enabled
	if hub != nil {
		hub.EnableDebug(debug)
	}
}
