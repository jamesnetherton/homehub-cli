package homehub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type apiTest struct {
	method          string
	methodArgs      []interface{}
	apiStubResponse string
	expectedResult  interface{}
	t               *testing.T
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func santizeString(target *string, regex string, replacement string) {
	re := regexp.MustCompile(regex)
	for _, match := range re.FindAllString(*target, -1) {
		*target = strings.Replace(*target, match, replacement, -1)
	}
}

func stubbedResponseHTTPHandler(apiStubResponse string, w http.ResponseWriter, r *http.Request) {
	var stubDataFile string
	if strings.HasSuffix(r.RequestURI, "/eventLog") {
		stubDataFile = "testdata/eventLog.txt"
	} else if strings.HasSuffix(r.RequestURI, "/stats.csv") {
		stubDataFile = "testdata/stats.csv"
	} else {
		stubDataFile = "testdata/" + apiStubResponse + "_response.json"
	}

	bytesRead, err := ioutil.ReadFile(stubDataFile)
	if err == nil {
		fmt.Fprintln(w, string(bytesRead))
	} else {
		fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": \"code\": 99999999, \"description\": \""+err.Error()+"\" }}")
	}
}

func proxiedResponseHTTPHandler(apiStubResponse string, url string, w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(r.Method, url+r.RequestURI, r.Body)
	req.ContentLength = r.ContentLength
	req.Form = r.Form
	req.Header = r.Header

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}

	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": { \"code\": 99999999, \"description\": \""+err.Error()+"\" }}}")
		return
	}

	defer httpResponse.Body.Close()
	bodyBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": { \"code\": 99999999, \"description\": \"Error reading proxied response\" }}}")
		return
	}

	body := string(bodyBytes[:])

	// Clean up MAC addresses
	santizeString(&body, "\\b([0-9a-fA-F]{2}:??){5}([0-9a-fA-F]{2})\\b", "11:AA:2B:33:44:5C")
	// Clean up IP addresses
	santizeString(&body, "\\b((25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)(\\.(25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)){3})\\b", "192.168.1.68")
	// Clean up timestampts
	santizeString(&body, "\\b([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})\\+([0-9]{4})\\b", "2016-08-30T19:48:55+0100")
	// Clean up serial number
	santizeString(&body, "\\b([0-9]{6})\\+([A-Z]{2})([0-9]{8})\\b", "123456+NQ98765432")

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": { \"code\": 99999999, \"description\": \"Error unmarshalling JSON response\" }}}")
		return
	}

	json, err := json.MarshalIndent(dat, "", "  ")
	if err != nil {
		fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": { \"code\": 99999999, \"description\": \"Error marshalling JSON response\" }}}")
		return
	}

	ioutil.WriteFile("testdata/"+apiStubResponse+"_response.json", json, 0644)
	fmt.Fprintln(w, body)
}

func mockAPIClientServer(apiStubResponse string) (*httptest.Server, *Hub) {
	defaultUsername := "admin"
	defaultPassword := "passw0rd"
	username := getEnv("HUB_USERNAME", defaultUsername)
	password := getEnv("HUB_PASSWORD", defaultPassword)
	debug := getEnv("HUB_DEBUG", "false")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username == defaultUsername && password == defaultPassword {
			stubbedResponseHTTPHandler(apiStubResponse, w, r)
		} else {
			proxiedResponseHTTPHandler(apiStubResponse, os.Getenv("HUB_URL"), w, r)
		}
	}))

	url := getEnv("HUB_URL", server.URL)
	hub := New(server.URL, username, password)

	if debug == "true" {
		hub.EnableDebug(true)
	}

	if url != server.URL {
		hub.Login()
	} else {
		hub.client.authData.userName = "admin"
		hub.client.authData.password = "admin"
		hub.client.authData.sessionID = "987879"
		hub.client.authData.nonce = "2355345"
	}

	return server, hub
}

func testAPIResponse(a *apiTest) {
	server, hub := mockAPIClientServer(a.apiStubResponse)
	defer server.Close()

	v := reflect.TypeOf(hub)

	apiMethod, _ := v.MethodByName(a.method)

	inputs := make([]reflect.Value, len(a.methodArgs)+1)
	for i := range a.methodArgs {
		inputs[i+1] = reflect.ValueOf(a.methodArgs[i])
	}

	inputs[0] = reflect.ValueOf(hub)
	resp := apiMethod.Func.Call(inputs)
	var result interface{}

	if resp[0].Type().String() == "string" {
		result = resp[0].String()
	} else if resp[0].Type().String() == "int" {
		result = int(resp[0].Int())
	} else if resp[0].Type().String() == "int64" {
		result = int64(resp[0].Int())
	} else if resp[0].Type().String() == "bool" {
		result = resp[0].Bool()
	} else if resp[0].Type().String() == "error" {
		if !resp[0].IsNil() {
			a.t.Fatalf("API method %s returned an unexpected error", a.method)
		}
	}

	if len(resp) > 1 {
		if !resp[1].IsNil() {
			if resp[1].Type().String() == "error" {
				result = fmt.Sprintf("%s", resp[1].Interface())
			}
		}
	}

	if result != a.expectedResult {
		a.t.Fatalf("API method %v returned '%v'. Expected '%v'", a.method, result, a.expectedResult)
	}
}

func TestBandwidthMonitor(t *testing.T) {
	server, hub := mockAPIClientServer("bandwidth_monitor")
	defer server.Close()

	res, err := hub.BandwidthMonitor()

	if err != nil {
		t.Fatalf("Error returned from BandwidthMonitor %s", err.Error())
	}

	if len(res.Entries) != 2 {
		t.Fatalf("Expected 2 bandwidth log entries but got %d", len(res.Entries))
	}

	if res.Entries[0].MACAddress != "a0:b1:c2:d3:e4:f5" {
		t.Fatalf("Expected bandwidth log entry 1 timestamp a0:b1:c2:d3:e4:f5 but got %s", res.Entries[0].MACAddress)
	}

	if res.Entries[0].Date != "2016-12-30" {
		t.Fatalf("Expected bandwidth log entry 1 date 2016-12-30 but got %s", res.Entries[0].Date)
	}

	if res.Entries[0].DownloadMegabytes != 10959 {
		t.Fatalf("Expected bandwidth log entry 1 download megabytes 10959 but got %s", res.Entries[0].DownloadMegabytes)
	}

	if res.Entries[0].UploadMegabytes != 1301 {
		t.Fatalf("Expected bandwidth log entry 1 upload megabytes 1301 but got %s", res.Entries[0].UploadMegabytes)
	}

	if res.Entries[1].MACAddress != "a1:b9:c8:d7:e6:f5" {
		t.Fatalf("Expected bandwidth log entry 2 timestamp a1:b9:c8:d7:e6:f5 but got %s", res.Entries[1].MACAddress)
	}

	if res.Entries[1].Date != "2016-12-31" {
		t.Fatalf("Expected bandwidth log entry 2 date 2016-12-31 but got %s", res.Entries[1].Date)
	}

	if res.Entries[1].DownloadMegabytes != 218 {
		t.Fatalf("Expected bandwidth log entry 2 download megabytes 218 but got %s", res.Entries[1].DownloadMegabytes)
	}

	if res.Entries[1].UploadMegabytes != 30 {
		t.Fatalf("Expected bandwidth log entry 2 upload megabytes ,30 but got %s", res.Entries[1].UploadMegabytes)
	}
}

func TestBroadbandProductType(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "BroadbandProductType",
		apiStubResponse: "interface_type",
		expectedResult:  "BT Infinity",
		t:               t,
	})
}

func TestConnectedDevices(t *testing.T) {
	server, hub := mockAPIClientServer("connected_devices")
	defer server.Close()

	res, err := hub.ConnectedDevices()

	if err != nil {
		t.Fatalf("Error returned from ConnectedDevices %s", err.Error())
	}

	if len(res) != 2 {
		t.Fatalf("Expected %d connected devices but got %d", 2, len(res))
	}

	if res[0].HostName != "foo.bar" {
		t.Fatalf("Expected device 1 to have host name foo.bar but got %s", res[0].HostName)
	}

	if len(res[0].IPv4Addresses) != 1 {
		t.Fatalf("Expected device 1 to have %d IPV4 addresses but got %d", 1, len(res[0].IPv4Addresses))
	}

	if len(res[0].IPv6Addresses) != 0 {
		t.Fatalf("Expected device 1 to have %d IPV6 addresses but got %d", 0, len(res[0].IPv6Addresses))
	}

	if res[1].HostName != "foo.bar.cheese" {
		t.Fatalf("Expected device 2 to have host name foo.bar but got %s", res[1].HostName)
	}

	if len(res[1].IPv4Addresses) != 1 {
		t.Fatalf("Expected device 2 to have %d IPV4 addresses but got %d", 1, len(res[1].IPv4Addresses))
	}

	if len(res[1].IPv6Addresses) != 0 {
		t.Fatalf("Expected device 2 to have %d IPV6 addresses but got %d", 0, len(res[1].IPv6Addresses))
	}
}

func TestDataPumpVersion(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DataPumpVersion",
		apiStubResponse: "data_pump_version",
		expectedResult:  "AfH042f.d26k1\n",
		t:               t,
	})
}

func TestDataReceived(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DataReceived",
		apiStubResponse: "data_received",
		expectedResult:  int64(99887766),
		t:               t,
	})
}

func TestDataSent(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DataSent",
		apiStubResponse: "data_sent",
		expectedResult:  int64(11223344),
		t:               t,
	})
}

func TestDeviceInfo(t *testing.T) {
	server, hub := mockAPIClientServer("device_info")
	defer server.Close()

	res, err := hub.DeviceInfo(2)

	if err != nil {
		t.Fatalf("Error returned from DeviceInfo %s", err.Error())
	}

	if res.HostName != "foo.bar" {
		t.Fatalf("Expected device to have host name foo.bar but got %s", res.HostName)
	}

	if len(res.IPv4Addresses) != 1 {
		t.Fatalf("Expected device to have %d IPV4 addresses but got %d", 1, len(res.IPv4Addresses))
	}

	if len(res.IPv6Addresses) != 0 {
		t.Fatalf("Expected device to have %d IPV6 addresses but got %d", 0, len(res.IPv6Addresses))
	}
}

func TestDhcpAuthoritative(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DhcpAuthoritative",
		apiStubResponse: "dhcp_authoritative",
		expectedResult:  true,
		t:               t,
	})
}

func TestDhcpPoolStart(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DhcpPoolStart",
		apiStubResponse: "dhcp_ipv4_pool_start",
		expectedResult:  "192.168.1.68",
		t:               t})
}

func TestDhcpPoolEnd(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DhcpPoolEnd",
		apiStubResponse: "dhcp_ipv4_pool_end",
		expectedResult:  "192.168.1.253",
		t:               t,
	})
}

func TestDhcpSubnetMask(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DhcpSubnetMask",
		apiStubResponse: "dhcp_subnet_mask",
		expectedResult:  "255.255.255.0",
		t:               t,
	})
}

func TestDownstreamSyncSpeed(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "DownstreamSyncSpeed",
		apiStubResponse: "downstream_curr_rate",
		expectedResult:  97543,
		t:               t,
	})
}

func TestEventLog(t *testing.T) {
	server, hub := mockAPIClientServer("event_log")
	defer server.Close()

	res, err := hub.EventLog()

	if err != nil {
		t.Fatalf("Error returned from EventLog %s", err.Error())
	}

	if len(res.Entries) != 2 {
		t.Fatalf("Expected 2 log entries but got %d", len(res.Entries))
	}

	if res.Entries[0].Timestamp != "01.03.2017 01:11:11" {
		t.Fatalf("Expected log entry 1 timestamp 01.03.2017 01:11:11 but got %s", res.Entries[0].Timestamp)
	}

	if res.Entries[0].Type != "INF" {
		t.Fatalf("Expected log entry 1 type INF but got %s", res.Entries[0].Type)
	}

	if res.Entries[0].Category != "WIFI" {
		t.Fatalf("Expected category entry 1 type WIFI but got %s", res.Entries[0].Category)
	}

	if res.Entries[0].Message != "Test log message 1" {
		t.Fatalf("Expected category entry 1 message 'Test log message 1' but got %s", res.Entries[0].Message)
	}

	if res.Entries[1].Timestamp != "02.03.2017 02:22:22" {
		t.Fatalf("Expected log entry 2 timestamp 02.03.2017 02:22:22 but got %s", res.Entries[1].Timestamp)
	}

	if res.Entries[1].Type != "WRN" {
		t.Fatalf("Expected log entry 2 type WRN but got %s", res.Entries[1].Type)
	}

	if res.Entries[1].Category != "TR69" {
		t.Fatalf("Expected category entry 2 type TR69 but got %s", res.Entries[1].Category)
	}

	if res.Entries[1].Message != "ppp1:TR69 ConnectionRequest: processing request from ACS" {
		t.Fatalf("Expected category entry 2 message 'ppp1:TR69 ConnectionRequest: processing request from ACS' but got %s", res.Entries[1].Message)
	}
}

func TestHardwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "HardwareVersion",
		apiStubResponse: "hardware_version",
		expectedResult:  "1.0",
		t:               t,
	})
}

func TestInternetConnectionStatus(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "InternetConnectionStatus",
		apiStubResponse: "wan_internet_status",
		expectedResult:  "UP",
		t:               t,
	})
}

func TestLightBrightness(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "LightBrightness",
		apiStubResponse: "hub_light_brightness",
		expectedResult:  50,
		t:               t,
	})
}

func TestLightBrightnessSet(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "LightBrightnessSet",
		methodArgs:      []interface{}{50},
		apiStubResponse: "hub_light_brightness_set",
		expectedResult:  nil,
		t:               t,
	})
}

func TestLightEnable(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "LightEnable",
		methodArgs:      []interface{}{true},
		apiStubResponse: "hub_light_enable",
		expectedResult:  nil,
		t:               t,
	})
}

func TestLightStatus(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "LightStatus",
		apiStubResponse: "hub_light_status",
		expectedResult:  "OFF",
		t:               t,
	})
}

func TestLoginSuccess(t *testing.T) {
	server, hub := mockAPIClientServer("login")
	defer server.Close()

	loggedIn, err := hub.Login()

	if err != nil {
		t.Error(err)
	}

	if !loggedIn {
		t.Errorf("Expected login to be successful")
	}
}

func TestLocalTime(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "LocalTime",
		apiStubResponse: "ntp_local_time",
		expectedResult:  "2016-08-30T19:48:55+0100",
		t:               t,
	})
}

func TestMaintenanceFirmwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "MaintenanceFirmwareVersion",
		apiStubResponse: "maintenance_firmware_version",
		expectedResult:  "SG0B000000AA",
		t:               t,
	})
}

func TestNatRules(t *testing.T) {
	server, hub := mockAPIClientServer("nat_rules")
	defer server.Close()

	res, err := hub.NatRules()

	if err != nil {
		t.Fatalf("Error returned from NatRules %s", err.Error())
	}

	if len(res) != 1 {
		t.Fatalf("Expected 1 NAT rule but got %d", len(res))
	}

	if res[0].Alias != "awesome-nat-rule-alias" {
		t.Fatalf("Expected NAT rule alias awesome-nat-rule-alias but got %s", res[0].Alias)
	}

	if res[0].AllExternalInterfaces != false {
		t.Fatalf("Expected NAT rule AllExternalInterfaces false but got %s", res[0].AllExternalInterfaces)
	}

	if res[0].Creator != "HUB_TESTER" {
		t.Fatalf("Expected NAT rule creator HUB_TESTER but got %s", res[0].Creator)
	}

	if res[0].Description != "Test NAT Rule description" {
		t.Fatalf("Expected NAT rule description Test NAT Rule description but got %s", res[0].Description)
	}

	if res[0].Enable != true {
		t.Fatalf("Expected NAT rule enable true but got %s", res[0].Enable)
	}

	if res[0].ExternalPort != 1111 {
		t.Fatalf("Expected NAT rule external port 1111 but got %d", res[0].ExternalPort)
	}

	if res[0].ExternalPortEndRange != 0 {
		t.Fatalf("Expected NAT rule external port end range 0 but got %d", res[0].ExternalPortEndRange)
	}

	if res[0].InternalClient != "192.168.1.68" {
		t.Fatalf("Expected NAT rule client IP 192.168.1.68 but got %s", res[0].InternalClient)
	}

	if res[0].InternalPort != 2222 {
		t.Fatalf("Expected NAT rule internal port 2222 but got %d", res[0].InternalPort)
	}

	if res[0].LeaseDuration != 60 {
		t.Fatalf("Expected NAT rule lease duration 60 but got %d", res[0].LeaseDuration)
	}

	if res[0].LeaseStart != "2016-08-30T19:48:55+0100" {
		t.Fatalf("Expected NAT rule lease start 2016-08-30T19:48:55+0100 but got %s", res[0].LeaseStart)
	}

	if res[0].Protocol != "TCP" {
		t.Fatalf("Expected NAT rule protocol TCP but got %s", res[0].Protocol)
	}

	if res[0].RemoteHost != "192.168.1.68" {
		t.Fatalf("Expected NAT rule remote host ip 192.168.1.68 but got %s", res[0].RemoteHost)
	}

	if res[0].Service != "TEST_SERVICE" {
		t.Fatalf("Expected NAT rule servie TEST_SERVICE but got %s", res[0].Service)
	}

	if res[0].Target != "ACCEPT" {
		t.Fatalf("Expected NAT rule target ACCEPT but got %s", res[0].Target)
	}

	if res[0].UID != 1 {
		t.Fatalf("Expected NAT rule type UID 1 but got %d", res[0].UID)
	}
}

func TestNatRule(t *testing.T) {
	server, hub := mockAPIClientServer("nat_rule")
	defer server.Close()

	natRule, err := hub.NatRule(1)

	if err != nil {
		t.Fatalf("Error returned from NatRule %s", err.Error())
	}

	if natRule.Alias != "awesome-nat-rule-alias" {
		t.Fatalf("Expected NAT rule alias awesome-nat-rule-alias but got %s", natRule.Alias)
	}

	if natRule.AllExternalInterfaces != false {
		t.Fatalf("Expected NAT rule AllExternalInterfaces false but got %s", natRule.AllExternalInterfaces)
	}

	if natRule.Creator != "HUB_TESTER" {
		t.Fatalf("Expected NAT rule creator HUB_TESTER but got %s", natRule.Creator)
	}

	if natRule.Description != "Test NAT Rule description" {
		t.Fatalf("Expected NAT rule description Test NAT Rule description but got %s", natRule.Description)
	}

	if natRule.Enable != true {
		t.Fatalf("Expected NAT rule enable true but got %s", natRule.Enable)
	}

	if natRule.ExternalPort != 1111 {
		t.Fatalf("Expected NAT rule external port 1111 but got %d", natRule.ExternalPort)
	}

	if natRule.ExternalPortEndRange != 0 {
		t.Fatalf("Expected NAT rule external port end range 0 but got %d", natRule.ExternalPortEndRange)
	}

	if natRule.InternalClient != "192.168.1.68" {
		t.Fatalf("Expected NAT rule client IP 192.168.1.68 but got %s", natRule.InternalClient)
	}

	if natRule.InternalPort != 2222 {
		t.Fatalf("Expected NAT rule internal port 2222 but got %d", natRule.InternalPort)
	}

	if natRule.LeaseDuration != 60 {
		t.Fatalf("Expected NAT rule lease duration 60 but got %d", natRule.LeaseDuration)
	}

	if natRule.LeaseStart != "2016-08-30T19:48:55+0100" {
		t.Fatalf("Expected NAT rule lease start 2016-08-30T19:48:55+0100 but got %s", natRule.LeaseStart)
	}

	if natRule.Protocol != "TCP" {
		t.Fatalf("Expected NAT rule protocol TCP but got %s", natRule.Protocol)
	}

	if natRule.RemoteHost != "192.168.1.68" {
		t.Fatalf("Expected NAT rule remote host ip 192.168.1.68 but got %s", natRule.RemoteHost)
	}

	if natRule.Service != "TEST_SERVICE" {
		t.Fatalf("Expected NAT rule servie TEST_SERVICE but got %s", natRule.Service)
	}

	if natRule.Target != "ACCEPT" {
		t.Fatalf("Expected NAT rule target ACCEPT but got %s", natRule.Target)
	}

	if natRule.UID != 1 {
		t.Fatalf("Expected NAT rule type UID 1 but got %d", natRule.UID)
	}
}

func TestNatRuleCreate(t *testing.T) {
	server, hub := mockAPIClientServer("nat_rule_create")
	defer server.Close()

	natRule := &NatRule{
		Enable:                false,
		Alias:                 "",
		ExternalInterface:     "",
		AllExternalInterfaces: false,
		LeaseDuration:         0,
		RemoteHost:            "",
		ExternalPort:          1111,
		ExternalPortEndRange:  1111,
		InternalInterface:     "",
		InternalPort:          0,
		Protocol:              "TCP",
		Service:               "Test Service",
		InternalClient:        "",
		Description:           "Test Description",
		Creator:               "JAMES",
		Target:                "REJECT",
		LeaseStart:            "",
	}

	err := hub.NatRuleCreate(natRule)

	if err != nil {
		t.Fatalf("Error returned from NatRuleCreate %s", err.Error())
	}

	if natRule.UID != 14 {
		t.Fatalf("Expected NAT rule UID 13 but was %d", natRule.UID)
	}
}

func TestNatRuleDelete(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "NatRuleDelete",
		methodArgs:      []interface{}{16},
		apiStubResponse: "nat_rule_delete",
		expectedResult:  nil,
		t:               t,
	})
}

func TestNatRuleUpdate(t *testing.T) {
	natRule := NatRule{
		UID:                   18,
		Enable:                true,
		Alias:                 "Updated Alias",
		ExternalInterface:     "",
		AllExternalInterfaces: false,
		LeaseDuration:         30,
		RemoteHost:            "",
		ExternalPort:          2222,
		ExternalPortEndRange:  2222,
		InternalInterface:     "",
		InternalPort:          0,
		Protocol:              "UDP",
		Service:               "FTP",
		InternalClient:        "",
		Description:           "Updated Test Description",
		Creator:               "HIDDEN",
		Target:                "DROP",
		LeaseStart:            "",
	}

	testAPIResponse(&apiTest{
		method:          "NatRuleUpdate",
		methodArgs:      []interface{}{natRule},
		apiStubResponse: "nat_rule_update",
		expectedResult:  nil,
		t:               t,
	})
}

func TestPublicIPAddress(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "PublicIPAddress",
		apiStubResponse: "public_ip4",
		expectedResult:  "192.168.1.68",
		t:               t,
	})
}

func TestPublicSubnetMask(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "PublicSubnetMask",
		apiStubResponse: "public_subnet_mask",
		expectedResult:  "255.255.255.255",
		t:               t,
	})
}

func TestReboot(t *testing.T) {
	// If we're testing against the real router, we don't want to reboot it midway through the test suite!
	if os.Getenv("HUB_USERNAME") == "" && os.Getenv("HUB_PASSWORD") == "" {
		testAPIResponse(&apiTest{
			method:          "Reboot",
			apiStubResponse: "reboot",
			expectedResult:  nil,
			t:               t,
		})
	}
}

func TestSambaHost(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "SambaHost",
		apiStubResponse: "samba_host",
		expectedResult:  "bthub,hub,bthomehub,api",
		t:               t,
	})
}

func TestSambaIP(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "SambaIP",
		apiStubResponse: "samba_ip",
		expectedResult:  "192.168.1.68",
		t:               t,
	})
}

func TestSerialNumber(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "SerialNumber",
		apiStubResponse: "serial_number",
		expectedResult:  "+123456+NQ98765432",
		t:               t,
	})
}

func TestSessionExpired(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "SerialNumber",
		apiStubResponse: "session_expired",
		expectedResult:  "Invalid user session",
		t:               t,
	})
}

func TestSoftwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "SoftwareVersion",
		apiStubResponse: "software_version",
		expectedResult:  "SG4B10002244",
		t:               t,
	})
}

func TestUpstreamSyncSpeed(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "UpstreamSyncSpeed",
		apiStubResponse: "upstream_curr_rate",
		expectedResult:  52121,
		t:               t,
	})
}

func TestVersion(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "Version",
		apiStubResponse: "hub_version",
		expectedResult:  "Home Hub 60 Type A",
		t:               t,
	})
}

func TestWiFiSecurityMode(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "WiFiSecurityMode",
		apiStubResponse: "wifi24_security_mode",
		expectedResult:  "ULTRA_SECURE_MODE",
		t:               t,
	})
}

func TestWiFiSSID(t *testing.T) {
	testAPIResponse(&apiTest{
		method:          "WiFiSSID",
		apiStubResponse: "wifi24_ssid",
		expectedResult:  "Click here for viruses",
		t:               t,
	})
}
