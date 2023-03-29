package main

import (
	"fmt"
	"github.com/deepch/vdk/format/mp4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// HTTPAPIServerStreamSaveToMP4 func
func HTTPAPIServerStreamSaveToMP4(c *gin.Context) {
	var err error

	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_save_mp4",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamSaveToMP4",
	})

	defer func() {
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "Close",
			}).Errorln(err)
		}
	}()

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("save", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamUnauthorized.Error())
		return
	}
	c.Writer.Write([]byte("await save started"))
	go func() {
		Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
		cid, ch, _, err := Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), MSE)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}

		defer Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "StreamCodecs",
			}).Errorln(err.Error())
			return
		}
		err = os.MkdirAll(fmt.Sprintf("save/%s/%s/", c.Param("uuid"), c.Param("channel")), 0755)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "MkdirAll",
			}).Errorln(err.Error())
		}
		f, err := os.Create(fmt.Sprintf("save/%s/%s/%s.mp4", c.Param("uuid"), c.Param("channel"), time.Now().String()))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "Create",
			}).Errorln(err.Error())
		}
		defer f.Close()

		muxer := mp4.NewMuxer(f)
		err = muxer.WriteHeader(codecs)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "WriteHeader",
			}).Errorln(err.Error())
			return
		}
		defer muxer.WriteTrailer()

		var videoStart bool
		controlExit := make(chan bool, 10)
		dur, err := time.ParseDuration(c.Param("duration"))
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "ParseDuration",
			}).Errorln(err.Error())
		}
		saveLimit := time.NewTimer(dur)
		noVideo := time.NewTimer(10 * time.Second)
		defer log.Println("client exit")
		for {
			select {
			case <-controlExit:
				requestLogger.WithFields(logrus.Fields{
					"call": "controlExit",
				}).Errorln("Client Reader Exit")
				return
			case <-saveLimit.C:
				requestLogger.WithFields(logrus.Fields{
					"call": "saveLimit",
				}).Errorln("Saved Limit End")
				return
			case <-noVideo.C:
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				if err = muxer.WritePacket(*pck); err != nil {
					return
				}
			}
		}
	}()
}
