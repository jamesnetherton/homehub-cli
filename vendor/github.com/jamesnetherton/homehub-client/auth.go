package homehub

type authData struct {
	url          string
	userName     string
	password     string
	sessionID    string
	nonce        string
	requestCount int
}
