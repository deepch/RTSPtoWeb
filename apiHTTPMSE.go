package main

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func HTTPAPIServerStreamMSE(ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "Close",
		}).Errorln(err)
	}()
	if !Storage.StreamChannelExist(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel"))) {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	Storage.StreamRun(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel")))
	err := ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "SetWriteDeadline",
		}).Errorln(err.Error())
		return
	}
	cid, ch, _, err := Storage.ClientAdd(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel")), MSE)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "ClientAdd",
		}).Errorln(err.Error())
		return
	}
	defer Storage.ClientDelete(ws.Request().FormValue("uuid"), cid, stringToInt(ws.Request().FormValue("channel")))
	codecs, err := Storage.StreamCodecs(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel")))
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "StreamCodecs",
		}).Errorln(err.Error())
		return
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "Send",
		}).Errorln(err.Error())
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "http_mse",
			"stream":  ws.Request().FormValue("uuid"),
			"channel": ws.Request().FormValue("channel"),
			"func":    "HTTPAPIServerStreamMSE",
			"call":    "Send",
		}).Errorln(err.Error())
		return
	}
	var videoStart bool
	go func() {
		defer func() {
			err := ws.Close()
			log.WithFields(logrus.Fields{
				"module":  "http_mse",
				"stream":  ws.Request().FormValue("uuid"),
				"channel": ws.Request().FormValue("channel"),
				"func":    "HTTPAPIServerStreamMSE",
				"call":    "Close",
			}).Errorln(err)
		}()
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "http_mse",
					"stream":  ws.Request().FormValue("uuid"),
					"channel": ws.Request().FormValue("channel"),
					"func":    "HTTPAPIServerStreamMSE",
					"call":    "Receive",
				}).Errorln(err.Error())
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-noVideo.C:
			log.WithFields(logrus.Fields{
				"module":  "http_mse",
				"stream":  ws.Request().FormValue("uuid"),
				"channel": ws.Request().FormValue("channel"),
				"func":    "HTTPAPIServerStreamMSE",
				"call":    "ErrorStreamNoVideo",
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
				log.WithFields(logrus.Fields{
					"module":  "http_mse",
					"stream":  ws.Request().FormValue("uuid"),
					"channel": ws.Request().FormValue("channel"),
					"func":    "HTTPAPIServerStreamMSE",
					"call":    "WritePacket",
				}).Errorln(err.Error())
				return
			}
			if ready {
				err := ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					log.WithFields(logrus.Fields{
						"module":  "http_mse",
						"stream":  ws.Request().FormValue("uuid"),
						"channel": ws.Request().FormValue("channel"),
						"func":    "HTTPAPIServerStreamMSE",
						"call":    "SetWriteDeadline",
					}).Errorln(err.Error())
					return
				}
				err = websocket.Message.Send(ws, buf)
				if err != nil {
					log.WithFields(logrus.Fields{
						"module":  "http_mse",
						"stream":  ws.Request().FormValue("uuid"),
						"channel": ws.Request().FormValue("channel"),
						"func":    "HTTPAPIServerStreamMSE",
						"call":    "Send",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}
}
