package homehub

import (
	"fmt"
	"time"
)

type bandwidthMonitorRequest struct {
	genericRequest
}

func newBandwidthMonitorRequest(authData *authData) (req *bandwidthMonitorRequest) {
	authData.requestCount++

	// TODO: Enable dates to be configurable
	now := time.Now()
	date := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	params := &parameters{
		StartDate: date,
		EndDate:   date,
	}

	a := action{
		ID:         0,
		Method:     methodUploadStatistics,
		XPath:      bandwidthMonitoring,
		Parameters: params,
	}

	var actions []action
	actions = append(actions, a)
	requestBody := newRequestBody(authData, actions)

	return &bandwidthMonitorRequest{
		genericRequest: genericRequest{
			*requestBody,
			*authData,
		},
	}
}
