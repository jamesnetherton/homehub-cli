package homehub

type client struct {
	authData authData
}

func newClient(URL string, username string, password string) *client {
	a := authData{
		url:      URL,
		userName: username,
		password: password,
	}
	return &client{a}
}

func (c *client) createRequest(method string, xpath string) (req *request) {
	return newRequest(&c.authData, method, xpath)
}

func (c *client) sendXPathRequest(xpath string) (result string, err error) {
	method := "getValue"

	// TODO: Clean up this piece of hackery
	if xpath == "Device" {
		method = "reboot"
	}

	req := c.createRequest(method, xpath)
	resp, err := req.send()

	if err == nil {
		return resp.getValue(), nil
	}

	return "", err
}
