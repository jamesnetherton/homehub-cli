package homehub

type firmwareSG4B1A struct {
	firmwareSG4B1
}

func (f *firmwareSG4B1A) downstreamSyncSpeedXPath() string {
	return deviceFastLinesLineTestParamsDownstreamCurrRate
}

func (f *firmwareSG4B1A) upstreamSyncSpeedXPath() string {
	return deviceFastLinesLineTestParamsUpstreamCurrRate
}
