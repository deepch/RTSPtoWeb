package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"
)

var debug bool

func NewStreamCore() *StorageST {
	argConfigPatch := flag.String("config", "config.json", "config patch (/etc/server/config.json or config.json)")
	argDebug := flag.Bool("debug", true, "set debug mode")
	debug = *argDebug
	flag.Parse()
	var tmp StorageST
	data, err := ioutil.ReadFile(*argConfigPatch)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	debug = tmp.Server.Debug
	for i, i2 := range tmp.Streams {
		i2.clients = make(map[string]ClientST)
		i2.ack = time.Now().Add(-255 * time.Hour)
		i2.hlsSegmentBuffer = make(map[int]Segment)
		tmp.Streams[i] = i2
	}
	return &tmp
}

/*
 Client Sections
*/

//ready
//ClientDelete Delete Client
func (obj *StorageST) SaveConfig() error {
	res, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", res, 0644)
	if err != nil {
		return err
	}
	return nil
}
