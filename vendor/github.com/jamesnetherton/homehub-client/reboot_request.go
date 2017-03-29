package homehub

type rebootRequest struct {
	genericRequest
}

func newRebootRequest(authData *authData) (req *rebootRequest) {
	authData.requestCount++

	params := &parameters{
		Source: "GUI",
	}

	a := action{
		ID:         0,
		Method:     methodReboot,
		XPath:      device,
		Parameters: params,
	}

	var actions []action
	actions = append(actions, a)
	requestBody := newRequestBody(authData, actions)

	return &rebootRequest{
		genericRequest: genericRequest{
			*requestBody,
			*authData,
		},
	}
}
