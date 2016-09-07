package homehub

type authData struct {
	url          string
	userName     string
	password     string
	sessionID    string
	nonce        string
	requestCount int
}

func (a *authData) isAuthenticated() bool {
	return a.url != "" && a.userName != "" && a.password != "" && a.sessionID != "" && a.nonce != ""
}
