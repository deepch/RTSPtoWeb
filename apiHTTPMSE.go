package main

import (
	"time"

	"github.com/gobwas/ws/wsutil"

	"github.com/gobwas/ws"

	"github.com/gin-gonic/gin"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreamMSE func
func HTTPAPIServerStreamMSE(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		return
	}

	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_mse",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamMSE",
	})

	defer func() {
		err = conn.Close()
		requestLogger.WithFields(logrus.Fields{
			"call": "Close",
		}).Errorln(err)
		log.Println("Client Full Exit")
	}()
	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("WS", c.Param("uuid"), c.Param("channel"), c.Param("token"), c.ClientIP()) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamUnauthorized.Error())
		return
	}

	Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "SetWriteDeadline",
		}).Errorln(err.Error())
		return
	}
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
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Send",
		}).Errorln(err.Error())
		return
	}
	var videoStart bool
	controlExit := make(chan bool, 10)
	noClient := time.NewTimer(10 * time.Second)
	go func() {
		defer func() {
			controlExit <- true
		}()
		for {
			header, _, err := wsutil.NextReader(conn, ws.StateServerSide)
			if err != nil {
				requestLogger.WithFields(logrus.Fields{
					"call": "Receive",
				}).Errorln(err.Error())
				return
			}
			switch header.OpCode {
			case ws.OpPong:
				noClient.Reset(10 * time.Second)
			case ws.OpClose:
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	pingTicker := time.NewTicker(500 * time.Millisecond)
	defer pingTicker.Stop()
	defer log.Println("client exit")
	for {
		select {

		case <-pingTicker.C:
			err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				return
			}
			buf, err := ws.CompileFrame(ws.NewPingFrame(nil))
			if err != nil {
				return
			}
			_, err = conn.Write(buf)
			if err != nil {
				return
			}
		case <-controlExit:
			requestLogger.WithFields(logrus.Fields{
				"call": "controlExit",
			}).Errorln("Client Reader Exit")
			return
		case <-noClient.C:
			requestLogger.WithFields(logrus.Fields{
				"call": "ErrorClientOffline",
			}).Errorln("Client OffLine Exit")
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
				err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "SetWriteDeadline",
					}).Errorln(err.Error())
					return
				}
				//err = websocket.Message.Send(ws, buf)
				err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf)
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
