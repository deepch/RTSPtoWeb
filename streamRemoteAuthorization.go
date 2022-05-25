package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type AuthorizationReq struct {
	Proto   string `json:"proto,omitempty"`
	Stream  string `json:"stream,omitempty"`
	Channel string `json:"channel,omitempty"`
	Token   string `json:"token,omitempty"`
	IP      string `json:"ip,omitempty"`
}

type AuthorizationRes struct {
	Status string `json:"status,omitempty"`
}

func RemoteAuthorization(proto string, stream string, channel string, token string, ip string) bool {

	if !Storage.ServerTokenEnable() {
		return true
	}

	response, err := http.PostForm(Storage.ServerTokenBackend(), url.Values{"proto": {proto}, "stream": {stream}, "channel": {channel}, "token": {token}, "ip":{ip}})

	if err != nil {
		return false
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	// debug
	// Convert the body to type string
	// sb := string(bodyBytes)
	// log.Printf(sb)

	if err != nil {
		return false
	}

	var res AuthorizationRes

	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		return false
	}

	if res.Status == "1" {
		return true
	}

	return false
}
