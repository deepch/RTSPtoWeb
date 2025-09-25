package main

import (
	"time"

	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	safeContext := c.Copy()
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_webrtc",
		"stream":  safeContext.Param("uuid"),
		"channel": safeContext.Param("channel"),
		"func":    "HTTPAPIServerStreamWebRTC",
	})

	if !Storage.StreamChannelExist(safeContext.Param("uuid"), safeContext.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("WebRTC", safeContext.Param("uuid"), safeContext.Param("channel"), safeContext.Query("token"), safeContext.ClientIP()) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamUnauthorized.Error())
		return
	}

	Storage.StreamChannelRun(safeContext.Param("uuid"), safeContext.Param("channel"))
	codecs, err := Storage.StreamChannelCodecs(safeContext.Param("uuid"), safeContext.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return
	}

	options := webrtc.Options{
	    ICEServers:    Storage.ServerICEServers(),
	    ICEUsername:   Storage.ServerICEUsername(),
	    ICECredential: Storage.ServerICECredential(),
	    PortMin:       Storage.ServerWebRTCPortMin(),
	    PortMax:       Storage.ServerWebRTCPortMax(),
	}
	
	//Ensures that ice_candidates is optional
	if len(Storage.ServerICECandidates()) > 0 {
	    options.ICECandidates = Storage.ServerICECandidates()
	}

	muxerWebRTC := webrtc.NewMuxer(options)

	answer, err := muxerWebRTC.WriteHeader(codecs, c.PostForm("data"))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}

	go func() {
		cid, ch, _, err := Storage.ClientAdd(safeContext.Param("uuid"), safeContext.Param("channel"), WEBRTC)
		if err != nil {
			c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}
		defer Storage.ClientDelete(safeContext.Param("uuid"), cid, safeContext.Param("channel"))
		defer muxerWebRTC.Close() // Close the WebRTC session when done
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
				return
			case pck, ok := <-ch:
				if !ok {
					// Channel closed, likely due to camera disconnection
					c.IndentedJSON(500, Message{Status: 0, Payload: "Camera disconnected"})
					requestLogger.WithFields(logrus.Fields{
						"call": "CameraDisconnected",
					}).Errorln("Camera disconnected")
					return
				}

				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "WritePacket",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}()
}
