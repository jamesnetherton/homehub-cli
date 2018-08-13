package homehub

import (
	"reflect"
	"strconv"
	"strings"
)

// NatRule represents a Home Hub NAT port forwarding rule
type NatRule struct {
	UID                   int    `json:"uid"`
	Enable                bool   `json:"Enable"`
	Alias                 string `json:"Alias"`
	ExternalInterface     string `json:"ExternalInterface"`
	AllExternalInterfaces bool   `json:"AllExternalInterfaces"`
	LeaseDuration         int    `json:"LeaseDuration"`
	RemoteHost            string `json:"RemoteHost"`
	ExternalPort          int    `json:"ExternalPort"`
	ExternalPortEndRange  int    `json:"ExternalPortEndRange"`
	InternalInterface     string `json:"InternalInterface"`
	InternalPort          int    `json:"InternalPort"`
	Protocol              string `json:"Protocol"`
	Service               string `json:"Service"`
	InternalClient        string `json:"InternalClient"`
	Description           string `json:"Description"`
	Creator               string `json:"Creator"`
	Target                string `json:"Target"`
	LeaseStart            string `json:"LeaseStart"`
}

type portMapping struct {
	NatRule `json:"PortMapping,omitempty"`
}

func (n *NatRule) getUpdateActions(xpath string) []action {
	var actions []action
	r := reflect.TypeOf(n).Elem()
	v := reflect.ValueOf(n).Elem()
	uid := v.FieldByName("UID").Int()

	id := 0
	for i := 0; i < v.NumField(); i++ {
		if r.Field(i).Name != "UID" &&
			r.Field(i).Name != "Creator" &&
			r.Field(i).Name != "Target" &&
			r.Field(i).Name != "LeaseStart" {
			action := action{
				ID:     id,
				Method: methodSetValue,
				XPath:  strings.Replace(xpath, "#", strconv.Itoa(int(uid)), 1) + "/" + r.Field(i).Name,
				Parameters: &parameters{
					Value: v.Field(i).Interface(),
				},
			}
			actions = append(actions, action)
			id++
		}
	}

	return actions
}
