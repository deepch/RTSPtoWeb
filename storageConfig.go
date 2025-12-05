package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/imdario/mergo"

	"github.com/liip/sheriff"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

// Command line flag global variables
var debug bool
var configFile string

//NewStreamCore do load config file
func NewStreamCore() *StorageST {
	flag.BoolVar(&debug, "debug", true, "set debug mode")
	flag.StringVar(&configFile, "config", "config.json", "config path (/etc/server/config.json or config.json)")
	flag.Parse()

	// Initialize security subsystem
	InitSecurity()

	var tmp StorageST

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "NewStreamCore",
			"call":   "godotenv.Load",
		}).Debugln("Error loading .env file (optional), using existing environment variables")
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "NewStreamCore",
			"call":   "ReadFile",
		}).Errorln(err.Error())
		os.Exit(1)
	}

	// Expand environment variables
	expandedData := []byte(os.ExpandEnv(string(data)))

	err = json.Unmarshal(expandedData, &tmp)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "NewStreamCore",
			"call":   "Unmarshal",
		}).Errorln(err.Error())
		os.Exit(1)
	}

	// Helper to decrypt channels
	decryptChannels := func(channels map[string]ChannelST) {
		for id, ch := range channels {
			if strings.HasPrefix(ch.URL, "enc:") {
				decrypted, err := Decrypt(ch.URL)
				if err != nil {
					log.WithFields(logrus.Fields{
						"module": "config",
						"func":   "NewStreamCore",
						"stream": id,
					}).Errorln("Failed to decrypt URL:", err)
				} else {
					ch.URL = decrypted
					channels[id] = ch
				}
			}
		}
	}

	if tmp.Server.Users == nil {
		tmp.Server.Users = make(map[string]UserST)
	}
	if tmp.Server.Sessions == nil {
		tmp.Server.Sessions = make(map[string]SessionST)
	}
	if envUsers := os.Getenv("RTSP_USERS"); envUsers != "" {
		users := strings.Split(envUsers, ",")
		for _, user := range users {
			parts := strings.Split(user, ":")
			if len(parts) == 3 {
				tmp.Server.Users[parts[0]] = UserST{Password: parts[1], Role: parts[2]}
			}
		}
	}

	// Decrypt global defaults
	// Note: We need to handle this if ChannelDefaults has a URL.
	// decryptChannels(map[string]ChannelST{"default": tmp.ChannelDefaults}) // Assuming ChannelDefaults matches structure but it's single Struct not map

	debug = tmp.Server.Debug
	for i, stream := range tmp.Streams {
		// Decrypt URLs in this stream's channels
		decryptChannels(stream.Channels)

		for i3, i4 := range stream.Channels {
			channel := tmp.ChannelDefaults
			err = mergo.Merge(&channel, i4)
			if err != nil {
				log.WithFields(logrus.Fields{
					"module": "config",
					"func":   "NewStreamCore",
					"call":   "Merge",
				}).Errorln(err.Error())
				os.Exit(1)
			}
			channel.clients = make(map[string]ClientST)
			channel.ack = time.Now().Add(-255 * time.Hour)
			channel.hlsSegmentBuffer = make(map[int]SegmentOld)
			channel.signals = make(chan int, 100)
			stream.Channels[i3] = channel
		}
		tmp.Streams[i] = stream
	}
	return &tmp
}

// ClientDelete Delete Client
func (obj *StorageST) SaveConfig() error {
	log.WithFields(logrus.Fields{
		"module": "config",
		"func":   "NewStreamCore",
	}).Debugln("Saving configuration to", configFile)
	v2, err := version.NewVersion("2.0.0")
	if err != nil {
		return err
	}

	// Create a deep copy or use Sheriff to filter, but we need to encrypt specific fields.
	// Sheriff marshals to map[string]interface{} (recurisvely).
	// It's easier to Marshal the object first, then modify the map, then write.
	// Or modify the object? Modifying the live object is dangerous if other threads read it.
	// BUT, obj.Streams is a map, writing to it is not thread safe if not locked.
	// The SaveConfig is likely called under lock or needs to be careful.
	// However, we can use Sheriff to marshal to interface{}, then traverse and encrypt.

	data, err := sheriff.Marshal(&sheriff.Options{
		Groups:     []string{"config"},
		ApiVersion: v2,
	}, obj)
	if err != nil {
		return err
	}

	// Encrypt URLs in the marshaled data
	// data is map[string]interface{}
	if streams, ok := data.(map[string]interface{})["streams"].(map[string]interface{}); ok {
		for _, streamVal := range streams {
			if stream, ok := streamVal.(map[string]interface{}); ok {
				if channels, ok := stream["channels"].(map[string]interface{}); ok {
					for _, chVal := range channels {
						if channel, ok := chVal.(map[string]interface{}); ok {
							if urlVal, ok := channel["url"].(string); ok && urlVal != "" {
								// Only encrypt if not env var (rudimentary check) and not already encrypted
								if !strings.HasPrefix(urlVal, "${") {
									encrypted, err := Encrypt(urlVal)
									if err == nil {
										channel["url"] = encrypted
									}
								}
							}
						}
					}
				}
			}
		}
	}

	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, res, 0644)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "config",
			"func":   "SaveConfig",
			"call":   "WriteFile",
		}).Errorln(err.Error())
		return err
	}
	return nil
}
