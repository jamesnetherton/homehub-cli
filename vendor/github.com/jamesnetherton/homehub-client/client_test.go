package homehub

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type apiTest struct {
	method          string
	apiStubResponse string
	expectedResult  string
	t               *testing.T
}

func mockAPIClientServer(apiStubResponse string) (*httptest.Server, *Hub) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytesRead, err := ioutil.ReadFile("testdata/" + apiStubResponse + "_response.json")
		if err == nil {
			fmt.Fprintln(w, string(bytesRead))
		} else {
			fmt.Fprintln(w, "{\"reply\": { \"uid\": 0 \"id\": 0 \"error\": \"code\": 99999999, \"description\": \"Error reading api stub response\" }}")
		}
	}))

	hub := New(server.URL, "admin", "passw0rd")
	return server, hub
}

func testAPIResponse(a *apiTest) {
	server, hub := mockAPIClientServer(a.apiStubResponse)
	defer server.Close()

	v := reflect.ValueOf(hub)

	// Simulate authentication before invoking target method
	hub.client.authData.userName = "admin"
	hub.client.authData.password = "admin"
	hub.client.authData.sessionID = "987879"
	hub.client.authData.nonce = "2355345"

	apiMethod := v.MethodByName(a.method)
	resp := apiMethod.Call(nil)

	result := resp[0].String()

	// Expect an empty response for a reboot
	if !resp[1].IsNil() && a.method != "Reboot" {
		a.t.Fatalf("API method %s failed", a.method)
	}

	if result != a.expectedResult {
		a.t.Fatalf("API method %s returned '%s'. Expected '%s'", a.method, result, a.expectedResult)
	}
}

func TestBroadbandProductType(t *testing.T) {
	testAPIResponse(&apiTest{"BroadbandProductType", "interface_type", "BT Infinity", t})
}

func TestDataPumpVersion(t *testing.T) {
	testAPIResponse(&apiTest{"DataPumpVersion", "data_pump_version", "AfH042f.d26k1\n", t})
}

func TestDataReceived(t *testing.T) {
	testAPIResponse(&apiTest{"DataReceived", "data_received", "99887766", t})
}

func TestDataSent(t *testing.T) {
	testAPIResponse(&apiTest{"DataSent", "data_sent", "11223344", t})
}

func TestDhcpPoolStart(t *testing.T) {
	testAPIResponse(&apiTest{"DhcpPoolStart", "dhcp_ipv4_pool_start", "192.168.1.64", t})
}

func TestDhcpPoolEnd(t *testing.T) {
	testAPIResponse(&apiTest{"DhcpPoolEnd", "dhcp_ipv4_pool_end", "192.168.1.253", t})
}

func TestDhcpSubnetMask(t *testing.T) {
	testAPIResponse(&apiTest{"DhcpSubnetMask", "dhcp_subnet_mask", "255.255.255.0", t})
}

func TestDownstreamSyncSpeed(t *testing.T) {
	testAPIResponse(&apiTest{"DownstreamSyncSpeed", "downstream_curr_rate", "97543", t})
}

func TestHardwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{"HardwareVersion", "hardware_version", "1.0", t})
}

func TestInternetConnectionStatus(t *testing.T) {
	testAPIResponse(&apiTest{"InternetConnectionStatus", "wan_internet_status", "UP", t})
}

func TestLightStatus(t *testing.T) {
	testAPIResponse(&apiTest{"LightStatus", "hub_light_status", "OFF", t})
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
	testAPIResponse(&apiTest{"LocalTime", "ntp_local_time", "2016-08-30T19:48:55+0100", t})
}

func TestMaintenanceFirmwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{"MaintenaceFirmwareVersion", "maintenance_firmware_version", "SG0B000000AA", t})
}

func TestPublicIPAddress(t *testing.T) {
	testAPIResponse(&apiTest{"PublicIPAddress", "public_ip4", "111.222.333.444", t})
}

func TestPublicSubnetMask(t *testing.T) {
	testAPIResponse(&apiTest{"PublicSubnetMask", "public_subnet_mask", "255.255.255.255", t})
}

func TestReboot(t *testing.T) {
	testAPIResponse(&apiTest{"Reboot", "reboot", "", t})
}

func TestSambaHost(t *testing.T) {
	testAPIResponse(&apiTest{"SambaHost", "samba_host", "bthub,hub,bthomehub,api", t})
}

func TestSambaIP(t *testing.T) {
	testAPIResponse(&apiTest{"SambaIP", "samba_ip", "192.168.1.254", t})
}

func TestSerialNumber(t *testing.T) {
	testAPIResponse(&apiTest{"SerialNumber", "serial_number", "+123456+NQ98765432", t})
}

func TestSoftwareVersion(t *testing.T) {
	testAPIResponse(&apiTest{"SoftwareVersion", "software_version", "SG4B100021AA", t})
}

func TestUpstreamSyncSpeed(t *testing.T) {
	testAPIResponse(&apiTest{"UpstreamSyncSpeed", "upstream_curr_rate", "52121", t})
}

func TestVersion(t *testing.T) {
	testAPIResponse(&apiTest{"Version", "hub_version", "Home Hub 60 Type A", t})
}

func TestWiFiSecurityMode(t *testing.T) {
	testAPIResponse(&apiTest{"WiFiSecurityMode", "wifi24_security_mode", "ULTRA_SECURE_MODE", t})
}

func TestWiFiSSID(t *testing.T) {
	testAPIResponse(&apiTest{"WiFiSSID", "wifi24_ssid", "Click here for viruses", t})
}
