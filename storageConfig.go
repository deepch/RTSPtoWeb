package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/liip/sheriff"

	"github.com/sirupsen/logrus"
)

//debug global
var debug bool

//NewStreamCore do load config file
func NewStreamCore() *StorageST {
	argConfigPatch := flag.String("config", "config.json", "config patch (/etc/server/config.json or config.json)")
	argDebug := flag.Bool("debug", true, "set debug mode")
	debug = *argDebug
	flag.Parse()
	var tmp StorageST
	data, err := ioutil.ReadFile(*argConfigPatch)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "NewStreamCore",
			"call":   "ReadFile",
		}).Errorln(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "NewStreamCore",
			"call":   "Unmarshal",
		}).Errorln(err.Error())
		os.Exit(1)
	}
	debug = tmp.Server.Debug
	for i, i2 := range tmp.Streams {
		for i3, i4 := range i2.Channels {
			i4.clients = make(map[string]ClientST)
			i4.ack = time.Now().Add(-255 * time.Hour)
			i4.hlsSegmentBuffer = make(map[int]SegmentOld)
			i4.signals = make(chan int, 100)
			i2.Channels[i3] = i4
		}
		tmp.Streams[i] = i2
	}
	return &tmp
}

//ClientDelete Delete Client
func (obj *StorageST) SaveConfig() error {
	v2, err := version.NewVersion("2.0.0")
	if err != nil {
		return err
	}
	data, err := sheriff.Marshal(&sheriff.Options{
		Groups:     []string{"config"},
		ApiVersion: v2,
	}, obj)
	if err != nil {
		return err
	}
	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", res, 0644)
	if err != nil {
		return err
	}
	return nil
}
