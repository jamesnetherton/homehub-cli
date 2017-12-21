package homehub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type loginRequest struct {
	genericRequest
}

func newLoginRequest(authData *authData) (req *loginRequest) {
	authData.requestCount = 0
	authData.sessionID = "0"

	newNss := newNss()
	var nssOptions []nss
	nssOptions = append(nssOptions, *newNss)

	contextFlags := &contextFlags{
		GetContentName: true,
		LocalTime:      true,
	}

	capabilityFlags := &capabilityFlags{
		Name:         true,
		DefaultValue: false,
		Restriction:  true,
		Description:  false,
	}

	sessionOptions := &sessionOptions{
		Nss:             nssOptions,
		Language:        "ident",
		ContextFlags:    *contextFlags,
		CapabilityDepth: 2,
		CapabilityFlags: *capabilityFlags,
		TimeFormat:      "ISO_8601",
	}

	parameters := &parameters{
		User:           authData.userName,
		Persistent:     "true",
		SessionOptions: sessionOptions,
	}

	a := action{
		ID:         0,
		Method:     methodLogin,
		Parameters: parameters,
	}

	var actions []action
	actions = append(actions, a)
	requestBody := newRequestBody(authData, actions)

	return &loginRequest{
		genericRequest: genericRequest{
			*requestBody,
			*authData,
		},
	}
}

func (r *loginRequest) send() (re *response, err error) {
	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	debug.Println(string(j))

	form := url.Values{}
	form.Add("req", string(j))

	sessionData := newSessionData(&r.authData)
	cj, _ := json.Marshal(sessionData)

	httpRequest, _ := http.NewRequest("POST", r.authData.url, strings.NewReader(form.Encode()))
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	httpRequest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")
	httpRequest.AddCookie(&http.Cookie{Name: "lang", Value: "en"})
	httpRequest.AddCookie(&http.Cookie{Name: "session", Value: url.QueryEscape(string(cj))})

	dump, _ := httputil.DumpRequest(httpRequest, true)
	debug.Println(string(dump))

	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	dump, _ = httputil.DumpResponse(httpResponse, true)
	debug.Println(string(dump))

	if httpResponse.StatusCode >= 400 {
		return nil, fmt.Errorf("Error processing request. Hub returned HTTP response code: %d", httpResponse.StatusCode)
	}

	defer httpResponse.Body.Close()
	response := &response{}
	bodyBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	response.body = string(bodyBytes[:])

	var responseBody = &responseBody{}
	json.Unmarshal(bodyBytes, responseBody)
	response.ResponseBody = *responseBody

	if responseBody.Reply != nil && responseBody.Reply.ReplyError.Description != "Ok" {
		err := errors.New(responseBody.Reply.ReplyError.Description)
		return nil, err
	}

	return response, nil
}
