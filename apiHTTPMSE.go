package main

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

//HTTPAPIServerStreamMSE func
func HTTPAPIServerStreamMSE(ws *websocket.Conn) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_mse",
		"stream":  ws.Request().FormValue("uuid"),
		"channel": ws.Request().FormValue("channel"),
		"func":    "HTTPAPIServerStreamMSE",
	})

	defer func() {
		err := ws.Close()
		requestLogger.WithFields(logrus.Fields{
			"call": "Close",
		}).Errorln(err)
		log.Println("Client Full Exit")
	}()
	if !Storage.StreamChannelExist(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel")) {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("WS", ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), ws.Request().FormValue("token"), ws.Request().RemoteAddr) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	Storage.StreamChannelRun(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"))
	err := ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "SetWriteDeadline",
		}).Errorln(err.Error())
		return
	}
	cid, ch, _, err := Storage.ClientAdd(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), MSE)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "ClientAdd",
		}).Errorln(err.Error())
		return
	}
	defer Storage.ClientDelete(ws.Request().FormValue("uuid"), cid, ws.Request().FormValue("channel"))
	codecs, err := Storage.StreamChannelCodecs(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}
	var videoStart bool
	controlExit := make(chan bool, 10)
	go func() {
		defer func() {
			controlExit <- true
		}()
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				requestLogger.WithFields(logrus.Fields{
					"call": "Receive",
				}).Errorln(err.Error())
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-controlExit:
			requestLogger.WithFields(logrus.Fields{
				"call": "controlExit",
			}).Errorln("Client Reader Exit")
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
			ready, buf, err := muxerMSE.WritePacket(*pck, false)
			if err != nil {
				requestLogger.WithFields(logrus.Fields{
					"call": "WritePacket",
				}).Errorln(err.Error())
				return
			}
			if ready {
				err := ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "SetWriteDeadline",
					}).Errorln(err.Error())
					return
				}
				err = websocket.Message.Send(ws, buf)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "Send",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}
}
