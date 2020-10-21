package main

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

//HTTPAPIServerStreamMSE func
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
		log.Println("Client Full Exit")
	}()
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 1")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 2")
	Storage.StreamChannelRun(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel")))
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 3")
	err := ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 4")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 5")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 6")
	defer Storage.ClientDelete(ws.Request().FormValue("uuid"), cid, stringToInt(ws.Request().FormValue("channel")))
	codecs, err := Storage.StreamChannelCodecs(ws.Request().FormValue("uuid"), stringToInt(ws.Request().FormValue("channel")))
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 7")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 8")
	muxerMSE := mp4f.NewMuxer(nil)
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 9")
	err = muxerMSE.WriteHeader(codecs)
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 10")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 11")
	meta, init := muxerMSE.GetInit(codecs)
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 12")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 13")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 14")
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
	log.Println(ws.Request().FormValue("uuid"), ws.Request().FormValue("channel"), "WS Step 15")
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-controlExit:
			log.WithFields(logrus.Fields{
				"module":  "http_mse",
				"stream":  ws.Request().FormValue("uuid"),
				"channel": ws.Request().FormValue("channel"),
				"func":    "HTTPAPIServerStreamMSE",
				"call":    "controlExit",
			}).Errorln("Client Reader Exit")
			return
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
