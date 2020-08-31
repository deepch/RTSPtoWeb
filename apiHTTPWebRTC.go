package main

import (
	"time"

	"github.com/deepch/vdk/format/webrtc"
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	if !Storage.StreamExist(c.Param("uuid")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		return
	}
	Storage.StreamRun(c.Param("uuid"))
	codecs, err := Storage.StreamCodecs(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	muxerWebRTC := webrtc.NewMuxer()
	answer, err := muxerWebRTC.WriteHeader(codecs, c.PostForm("data"))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	go func() {
		cid, ch, err := Storage.ClientAdd(c.Param("uuid"))
		if err != nil {
			c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
			loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
			return
		}
		defer Storage.ClientDelete(c.Param("uuid"), cid)
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
					return
				}
			}
		}
	}()
}
