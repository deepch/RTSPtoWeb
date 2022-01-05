package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
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

	buf, err := json.Marshal(&AuthorizationReq{proto, stream, channel, token, ip})

	if err != nil {
		return false
	}

	request, err := http.NewRequest("POST", Storage.ServerTokenBackend(), bytes.NewBuffer(buf))

	if err != nil {
		return false
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Do(request)

	if err != nil {
		return false
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

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
