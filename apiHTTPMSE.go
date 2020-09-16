package main

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"golang.org/x/net/websocket"
)

func HTTPAPIServerStreamMSE(ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err})
	}()
	if !Storage.StreamChannelExist(ws.Request().FormValue("uuid"), 0) {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		return
	}
	err := ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	cid, ch, err := Storage.ClientAdd(ws.Request().FormValue("uuid"), 0)
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	defer Storage.ClientDelete(ws.Request().FormValue("uuid"), cid, 0)
	Storage.StreamRun(ws.Request().FormValue("uuid"), 0)
	codecs, err := Storage.StreamCodecs(ws.Request().FormValue("uuid"), 0)
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	var videoStart bool
	go func() {
		defer func() {
			err := ws.Close()
			loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err})
		}()
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-noVideo.C:
			loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
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
				loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
				return
			}
			if ready {
				err := ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
					return
				}
				err = websocket.Message.Send(ws, buf)
				if err != nil {
					loggingPrintln(ws.Request().FormValue("uuid"), Message{Status: 0, Payload: err.Error()})
					return
				}
			}
		}
	}
}
