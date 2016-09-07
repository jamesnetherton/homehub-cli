package homehub

import (
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
			switch {
			case strings.Contains(params.Capability.Type, "uint32"):
				value = strconv.FormatFloat(vo.Float(), 'f', -1, 64)
				break
			default:
				value = vo.String()
			}
		}
	}

	return value
}
