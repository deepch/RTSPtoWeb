package main

import (
	"time"

	"github.com/deepch/vdk/format/webrtc"
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreamWebRTC need work
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	uuid := c.Param("uuid")
	data := c.PostForm("data")
	if !Storage.StreamExist(uuid) {
		c.IndentedJSON(500, ErrorNotFound)
		return
	}
	Storage.StreamRun(uuid)
	codecs, err := Storage.StreamCodecs(uuid)
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	muxerWebRTC := webrtc.NewMuxer()
	answer, err := muxerWebRTC.WriteHeader(codecs, data)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	go func() {
		cid, ch, err := Storage.ClientAdd(uuid)
		if err != nil {
			return
		}
		defer Storage.ClientDelete(uuid, cid)
		var videoStart bool
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
				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					return
				}
			}
		}
	}()
}
