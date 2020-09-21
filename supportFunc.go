package main

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

//Default streams signals
const (
	SignalStreamRestart = iota ///< Y   Restart
	SignalStreamStop
	SignalStreamClient
)

//generateUUID function make random uuid for clients and stream
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

//stringToInt convert string to int if err to zero
func stringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

//stringInBetween fin char to char sub string
func stringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	str = str[s+len(start):]
	e := strings.Index(str, end)
	if e == -1 {
		return
	}
	str = str[:e]
	return str
}
