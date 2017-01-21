package homehub

type xpathRequest struct {
	genericRequest
}

func newXPathRequest(authData *authData, xpath string, method string, value interface{}) (req *xpathRequest) {
	var (
		caps    *capabilityFlags
		options *interfaceOptions
		params  *parameters
	)

	authData.requestCount++

	if method == methodGetValue {
		caps = &capabilityFlags{
			Interface: true,
		}

		options = &interfaceOptions{
			CapabilityFlags: *caps,
		}
	} else {
		params = &parameters{
			Value: value,
		}
	}

	a := action{
		ID:               0,
		Method:           method,
		XPath:            xpath,
		InterfaceOptions: options,
		Parameters:       params,
	}

	var actions []action
	actions = append(actions, a)
	requestBody := newRequestBody(authData, actions)

	return &xpathRequest{
		genericRequest: genericRequest{
			*requestBody,
			*authData,
		},
	}
}
