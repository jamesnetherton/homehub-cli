package homehub

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type client struct {
	authData authData
}

func newClient(URL string, username string, password string) *client {
	auth := authData{
		url:      URL,
		userName: username,
		password: password,
	}

	return &client{auth}
}

func (c *client) doReboot() (err error) {
	req := newRebootRequest(&c.authData)
	_, err = req.send()
	return err
}

func (c *client) getBandwidthUsage() (result string, err error) {
	bandwidthMonitorRequest := newBandwidthMonitorRequest(&c.authData)
	req := newHubResourceRequest(&c.authData, c.authData.url, bandwidthMonitorRequest)
	resp, err := req.send()
	if err != nil {
		return "", err
	}

	return resp.body, nil
}

func (c *client) getEventLog() (result string, err error) {
	eventLogRequest := newEventLogRequest(&c.authData)
	req := newHubResourceRequest(&c.authData, c.authData.url, eventLogRequest)
	resp, err := req.send()

	if err != nil {
		return "", err
	}

	return resp.body, nil
}

func (c *client) getXPathValueString(xpath string) (result string, err error) {
	resp, err := c.doXPathRequest(xpath)

	if err == nil {
		if resp.ResponseBody.Reply != nil {
			params := resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
			vo := reflect.ValueOf(params.Value)

			if getTypeMapping(params.Capability.Type) != "string" {
				return "", errors.New("Expected response value to be of type string but was " + getTypeMapping(params.Capability.Type))
			}

			return vo.String(), nil
		}
	}

	return "", err
}

func (c *client) getXPathValueInt(xpath string) (result int, err error) {
	resp, err := c.doXPathRequest(xpath)

	if err == nil {
		if resp.ResponseBody.Reply != nil {
			params := resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
			vo := reflect.ValueOf(params.Value)
			if getTypeMapping(params.Capability.Type) != "int" {
				return -1, errors.New("Expected response value to be of type int but was " + getTypeMapping(params.Capability.Type))
			}
			if vo.Type().String() == "float64" {
				return int(vo.Float()), nil
			} else if vo.Type().String() == "string" {
				i, _ := strconv.Atoi(vo.String())
				return i, nil
			}
			return int(vo.Int()), nil
		}
	}

	return -1, err
}

func (c *client) getXPathValueInt64(xpath string) (result int64, err error) {
	resp, err := c.doXPathRequest(xpath)

	if err == nil {
		if resp.ResponseBody.Reply != nil {
			params := resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
			vo := reflect.ValueOf(params.Value)
			if getTypeMapping(params.Capability.Type) != "int64" {
				return -1, errors.New("Expected response value to be of type int64 but was " + getTypeMapping(params.Capability.Type))
			}
			if vo.Type().String() == "float64" {
				return int64(vo.Float()), nil
			} else if vo.Type().String() == "string" {
				i, _ := strconv.Atoi(vo.String())
				return int64(i), nil
			}
			return int64(vo.Int()), nil
		}
	}

	return -1, err
}

func (c *client) getXPathValueBool(xpath string) (result bool, err error) {
	resp, err := c.doXPathRequest(xpath)

	if err == nil {
		if resp.ResponseBody.Reply != nil {
			params := resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
			vo := reflect.ValueOf(params.Value)
			if getTypeMapping(params.Capability.Type) != "bool" {
				return false, errors.New("Expected response value to be of type bool but was " + getTypeMapping(params.Capability.Type))
			}
			if vo.Type().String() == "string" {
				b, _ := strconv.ParseBool(vo.String())
				return b, nil
			}
			return vo.Bool(), nil
		}
	}

	return false, err
}

func (c *client) getXPathValueType(xpath string, valueType reflect.Type) (result interface{}, err error) {
	req := newXPathRequest(&c.authData, xpath, methodGetValue, nil)
	resp, err := req.send()

	if err != nil {
		return nil, err
	}

	return getMapValue(resp, valueType), nil
}

func (c *client) getXPathValues(xpath string, valueType reflect.Type) (values []interface{}, err error) {
	req := newXPathRequest(&c.authData, xpath, methodGetValue, nil)
	resp, err := req.send()

	if err == nil {
		var values []interface{}

		if resp.ResponseBody.Reply != nil {
			for _, action := range resp.ResponseBody.Reply.ResponseActions {
				c := action.ResponseCallbacks[0]
				if c.XPath == xpath {
					p := c.Parameters
					if strings.HasPrefix(fmt.Sprintf("%s", p.Value), "[") {
						v := reflect.New(valueType).Interface()
						x, _ := json.Marshal(p.Value)
						json.Unmarshal(x, v)
						array := reflect.ValueOf(v).Elem()
						for i := 0; i < array.Len(); i++ {
							values = append(values, array.Index(i).Interface())
						}
					}
				}
			}
		}

		return values, nil
	}

	return nil, err
}

func (c *client) setXPathValue(xpath string, value interface{}) (err error) {
	req := newXPathRequest(&c.authData, xpath, methodSetValue, value)
	_, err = req.send()
	return err
}

func (c *client) setXPathValues(actions []action) (err error) {
	req := newMultiXPathRequest(&c.authData, actions)
	_, err = req.send()
	return err
}

func (c *client) addChildXPathValue(xpath string, value interface{}) (result int, err error) {
	req := newXPathRequest(&c.authData, xpath, methodAddChild, value)
	resp, err := req.send()

	if err == nil {
		responseXPath := resp.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].XPath

		re := regexp.MustCompile("'(.*)'")
		matches := re.FindAllStringSubmatch(responseXPath, -1)
		if len(matches) > 0 {
			return strconv.Atoi(matches[0][1])
		}
	}

	return -1, err
}

func (c *client) deleteChildXPathValue(xpath string) (err error) {
	req := newXPathRequest(&c.authData, xpath, methodDeleteValue, nil)
	_, err = req.send()
	return err
}

func (c *client) doXPathRequest(xpath string) (response *response, err error) {
	return newXPathRequest(&c.authData, xpath, methodGetValue, nil).send()
}

func getMapValue(response *response, valueType reflect.Type) (result interface{}) {
	var xPathValueType interface{}
	if response.ResponseBody.Reply != nil {
		params := response.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
		if strings.HasPrefix(fmt.Sprintf("%s", params.Value), "map[") {
			v := reflect.New(valueType).Interface()
			x, _ := json.Marshal(params.Value)
			json.Unmarshal(x, v)
			xPathValueType = reflect.ValueOf(v).Interface()
		}
	}
	return xPathValueType
}
