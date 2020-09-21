package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	//TODO: next add write to file
	if !debug {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(Storage.ServerLogLevel())
}
