package homehub

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type request interface {
	send() (re *response, err error)
}

type genericRequest struct {
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
	ID             int             `json:"id,omitempty"`
	Nonce          string          `json:"nonce,omitempty"`
	Persistent     string          `json:"persistent,omitempty"`
	SessionOptions *sessionOptions `json:"session-options,omitempty"`
	User           string          `json:"user,omitempty"`
	Value          interface{}     `json:"value,omitempty"`
	Capability     *capability     `json:"capability,omitempty"`
	URI            string          `json:"uri,omitempty"`
	Data           string          `json:"data,omitempty"`
	FileName       string          `json:"FileName,omitempty"`
	StartDate      string          `json:"startDate,omitempty"`
	EndDate        string          `json:"endDate,omitempty"`
	Source         string          `json:"source,omitempty"`
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

func (r *genericRequest) send() (re *response, err error) {
	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	debug.Println(string(j))

	form := url.Values{}
	form.Add("req", string(j))

	httpRequest, _ := http.NewRequest("POST", r.authData.url, strings.NewReader(form.Encode()))
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	httpRequest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")

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

func (r *genericRequest) isLogin() bool {
	return r.Body.Actions[0].Method == "logIn"
}

func newNss() *nss {
	return &nss{Name: "gtw", URI: "http://sagemcom.com/gateway-data"}
}

func newRequestBody(authData *authData, actions []action) *requestBody {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	cnonce := random.Intn(math.MaxInt32)

	var ha1 string
	if authData.nonce != "" {
		ha1 = hexmd5(authData.userName + ":" + authData.nonce + ":" + authData.password)
	} else {
		ha1 = hexmd5(authData.userName + "::" + authData.password)
	}
	authKey := hexmd5(ha1 + ":" + strconv.Itoa(authData.requestCount) + ":" + strconv.Itoa(cnonce) + ":JSON:/cgi/json-req")

	return &requestBody{
		ID:                authData.requestCount,
		SessionID:         authData.sessionID,
		SessionExpiryTime: "",
		Priority:          false,
		Actions:           actions,
		CNonce:            cnonce,
		AuthKey:           authKey,
	}
}

func newSessionData(authData *authData) *sessionData {
	newNss := newNss()
	var nssOptions []nss
	nssOptions = append(nssOptions, *newNss)

	dataModel := &dataModel{
		Name: "Internal",
		Nss:  nssOptions,
	}

	sessionID, _ := strconv.Atoi(authData.sessionID)
	authKey := hexmd5(authData.userName + ":" + authData.nonce + ":" + authData.password)
	ha1 := authKey[:10] + authData.password + authKey[10:len(authKey)]

	return &sessionData{
		ID:        authData.requestCount,
		SessionID: sessionID,
		Basic:     false,
		User:      authData.userName,
		DataModel: *dataModel,
		Ha1:       ha1,
		Nonce:     authData.nonce,
	}
}

func hexmd5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
