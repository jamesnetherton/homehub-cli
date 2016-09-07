package homehub

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type request struct {
	Body     requestBody `json:"request"`
	authData authData
}

type requestBody struct {
	ID                int      `json:"id"`
	SessionID         string   `json:"session-id"`
	SessionExpiryTime string   `json:"-"`
	Priority          bool     `json:"priority"`
	Actions           []action `json:"actions"`
	CNonce            int      `json:"cnonce"`
	AuthKey           string   `json:"auth-key"`
}

type action struct {
	ID               int               `json:"id"`
	Method           string            `json:"method"`
	XPath            string            `json:"xpath,omitempty"`
	Parameters       *parameters       `json:"parameters,omitempty"`
	InterfaceOptions *interfaceOptions `json:"options,omitempty"`
}

type parameters struct {
	ID             int            `json:"id,omitempty"`
	Nonce          string         `json:"nonce,omitempty"`
	Persistent     string         `json:"persistent,omitempty"`
	SessionOptions sessionOptions `json:"session-options,omitempty"`
	User           string         `json:"user,omitempty"`
	Value          interface{}    `json:"value,omitempty"`
	Capability     *capability    `json:"capability,omitempty"`
}

type interfaceOptions struct {
	CapabilityFlags capabilityFlags `json:"capability-flags"`
}

type sessionOptions struct {
	Nss             []nss           `json:"nss"`
	Language        string          `json:"language"`
	ContextFlags    contextFlags    `json:"context-flags"`
	CapabilityFlags capabilityFlags `json:"capability-flags"`
	CapabilityDepth int             `json:"capability-depth"`
	TimeFormat      string          `json:"time-format"`
}

type nss struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type contextFlags struct {
	GetContentName bool `json:"get-content-name"`
	LocalTime      bool `json:"local-time"`
}

type capabilityFlags struct {
	Name         bool `json:"name,omitempty"`
	DefaultValue bool `json:"default-value,omitempty"`
	Restriction  bool `json:"restriction,omitempty"`
	Description  bool `json:"description,omitempty"`
	Interface    bool `json:"interface,omitempty"`
}

type capability struct {
	Type string `json:"type"`
}

type sessionData struct {
	ID        int       `json:"req_id"`
	SessionID int       `json:"sess_id"`
	Basic     bool      `json:"basic"`
	User      string    `json:"user"`
	DataModel dataModel `json:"dataModel"`
	Ha1       string    `json:"ha1"`
	Nonce     string    `json:"nonce"`
}

type dataModel struct {
	Name string `json:"name"`
	Nss  []nss  `json:"nss"`
}

func newRequest(authData *authData, method string, xpath string) (req *request) {
	var a action
	var priority bool

	// TODO: This is horrible and needs tidying up
	if method == "logIn" {
		authData.requestCount = 0
		authData.sessionID = "0"
		priority = true

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
			SessionOptions: *sessionOptions,
		}

		a = action{
			ID:         0,
			Method:     method,
			Parameters: parameters,
		}
	} else {
		authData.requestCount++
		priority = false

		capabilityFlags := &capabilityFlags{
			Interface: true,
		}

		interfaceOptions := &interfaceOptions{
			CapabilityFlags: *capabilityFlags,
		}

		a = action{
			ID:               0,
			Method:           method,
			XPath:            xpath,
			InterfaceOptions: interfaceOptions,
		}
	}

	var actions []action
	actions = append(actions, a)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	cnonce := random.Intn(math.MaxInt32)

	var ha1 string
	if authData.nonce != "" {
		ha1 = hexmd5(authData.userName + ":" + authData.nonce + ":" + hexmd5(authData.password))
	} else {
		ha1 = hexmd5(authData.userName + "::" + hexmd5(authData.password))
	}
	authKey := hexmd5(ha1 + ":" + strconv.Itoa(authData.requestCount) + ":" + strconv.Itoa(cnonce) + ":JSON:/cgi/json-req")

	requestBody := &requestBody{
		ID:                authData.requestCount,
		SessionID:         authData.sessionID,
		SessionExpiryTime: "",
		Priority:          priority,
		Actions:           actions,
		CNonce:            cnonce,
		AuthKey:           authKey,
	}

	return &request{*requestBody, *authData}
}

func (r *request) send() (re *response, err error) {

	if !r.authData.isAuthenticated() && !r.isLogin() {
		return nil, errors.New("User not logged in")
	}

	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	log.Println(string(j))

	form := url.Values{}
	form.Add("req", string(j))

	httpRequest, _ := http.NewRequest("POST", r.authData.url, strings.NewReader(form.Encode()))

	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	httpRequest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")

	newNss := newNss()
	var nssOptions []nss
	nssOptions = append(nssOptions, *newNss)

	dataModel := &dataModel{
		Name: "Internal",
		Nss:  nssOptions,
	}

	if r.isLogin() {
		sessionID, _ := strconv.Atoi(r.authData.sessionID)
		authKey := hexmd5(r.authData.userName + ":" + r.authData.nonce + ":" + hexmd5(r.authData.password))
		ha1 := authKey[:10] + hexmd5(r.authData.password) + authKey[10:len(authKey)]

		sessionData := &sessionData{
			ID:        r.authData.requestCount,
			SessionID: sessionID,
			Basic:     false,
			User:      r.authData.userName,
			DataModel: *dataModel,
			Ha1:       ha1,
			Nonce:     r.authData.nonce,
		}

		cj, _ := json.Marshal(sessionData)

		httpRequest.AddCookie(&http.Cookie{Name: "lang", Value: "en"})
		httpRequest.AddCookie(&http.Cookie{Name: "session", Value: url.QueryEscape(string(cj))})
	}

	dump, _ := httputil.DumpRequest(httpRequest, true)
	log.Println(string(dump))

	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	dump, _ = httputil.DumpResponse(httpResponse, true)
	log.Println(string(dump))

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

	// TODO: This logic doesn't really belong here
	if responseBody.Reply != nil && responseBody.Reply.ReplyError.Description != "Ok" {
		err := errors.New(responseBody.Reply.ReplyError.Description)
		return nil, err
	}

	return response, nil
}

func (r *request) isLogin() bool {
	return r.Body.Actions[0].Method == "logIn"
}

func newNss() *nss {
	return &nss{Name: "gtw", URI: "http://sagemcom.com/gateway-data"}
}

func hexmd5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
