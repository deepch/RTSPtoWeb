package main

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"golang.org/x/net/websocket"
)

func HTTPAPIServerStreamMSE(ws *websocket.Conn) {
	defer ws.Close()
	uuid := ws.Request().FormValue("uuid")
	if !Storage.StreamExist(uuid) {
		return
	}
	err := ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return
	}
	cid, ch, err := Storage.ClientAdd(uuid)
	if err != nil {
		return
	}
	defer Storage.ClientDelete(uuid, cid)
	Storage.StreamRun(uuid)
	codecs, err := Storage.StreamCodecs(uuid)
	if err != nil {
		return
	}

	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		return
	}
	var videoStart bool
	go func() {
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				err = ws.Close()
				if err != nil {
					return
				}
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-noVideo.C:
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
				return
			}
			if ready {
				err := ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					return
				}
				err = websocket.Message.Send(ws, buf)
				if err != nil {
					return
				}
			}
		}
	}
}
