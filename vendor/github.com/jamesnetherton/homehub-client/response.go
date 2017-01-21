package homehub

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type response struct {
	body         string
	ResponseBody responseBody
}

type responseBody struct {
	Reply *reply `json:"reply"`
}

type reply struct {
	UID             int              `json:"uid"`
	ID              int              `json:"id"`
	ReplyError      replyError       `json:"error"`
	ResponseActions []responseAction `json:"actions"`
}

type replyError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type responseAction struct {
	UID               int                `json:"uid"`
	ID                int                `json:"id"`
	ReplyError        replyError         `json:"error"`
	ResponseCallbacks []responseCallback `json:"callbacks"`
	ResponseEvents    []responseEvent    `json:"events"`
}

type responseCallback struct {
	UID        int        `json:"uid"`
	Result     result     `json:"result"`
	XPath      string     `json:"xpath"`
	Parameters parameters `json:"parameters"`
}

type result struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type responseEvent struct {
	//TODO: Events not supported right now
}

func (r *response) getValue() string {
	value := ""

	if r.ResponseBody.Reply != nil {
		params := r.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
		vo := reflect.ValueOf(params.Value)

		if params.Capability == nil {
			value = vo.String()
		} else {
			capType := params.Capability.Type
			switch {
			case strings.Contains(capType, "deviceconfig:LastSuccesfulWanType"):
				value = vo.String()
				break
			case strings.Contains(capType, "int32"):
				value = strconv.FormatFloat(vo.Float(), 'f', -1, 64)
				break
			case strings.Contains(capType, "boolean"):
				value = strconv.FormatBool(vo.Bool())
				break
			default:
				value = vo.String()
			}
		}
	}

	return value
}

func (r *response) getValues(xpath string) [][]value {
	var res [][]value

	if r.ResponseBody.Reply != nil {
		for _, action := range r.ResponseBody.Reply.ResponseActions {
			c := action.ResponseCallbacks[0]
			if c.XPath == xpath {
				p := c.Parameters
				if strings.HasPrefix(fmt.Sprintf("%s", p.Value), "[") {
					v := &[]value{}
					x, _ := json.Marshal(p.Value)
					json.Unmarshal(x, v)
					res = append(res, *v)
				}
			}
		}
	}

	return res
}

func (r *response) getHost() *host {
	var h *host

	if r.ResponseBody.Reply != nil {
		params := r.ResponseBody.Reply.ResponseActions[0].ResponseCallbacks[0].Parameters
		if strings.HasPrefix(fmt.Sprintf("%s", params.Value), "map[Host") {
			h = &host{}
			x, _ := json.Marshal(params.Value)
			json.Unmarshal(x, h)
		}
	}

	return h
}
