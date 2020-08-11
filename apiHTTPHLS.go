package main

import (
	"bytes"
	"log"
	"time"

	"github.com/deepch/vdk/format/ts"
	"github.com/gin-gonic/gin"
)

//ready
//HTTPAPIServerStreamHLSTS send client m3u8 play list
func HTTPAPIServerStreamHLSM3U8(c *gin.Context) {
	uuid := c.Param("uuid")
	if !Storage.StreamExist(uuid) {
		c.IndentedJSON(500, "Stream Not Found")
		return
	}
	c.Header("Content-Type", "application/x-mpegURL")
	Storage.StreamRun(uuid)
	//If stream mode on_demand need wait ready segment's
	for i := 0; i < 40; i++ {
		index, seq, err := Storage.StreamHLSm3u8(uuid)
		if err != nil {
			c.IndentedJSON(500, err)
			return
		}
		if seq >= 6 {
			_, err := c.Writer.Write([]byte(index))
			if err != nil {
				c.IndentedJSON(400, err.Error())
				return
			}
			return
		}
		time.Sleep(1 * time.Second)
	}
}

//ready
//HTTPAPIServerStreamHLSTS send client ts segment
func HTTPAPIServerStreamHLSTS(c *gin.Context) {
	uuid := c.Param("uuid")
	//Check Has Stream
	if !Storage.StreamExist(uuid) {
		log.Println("Not Found Error")
		return
	}
	outfile := bytes.NewBuffer([]byte{})
	codecs, err := Storage.StreamCodecs(uuid)
	if err != nil {
		c.IndentedJSON(500, err.Error())
		return
	}
	Muxer := ts.NewMuxer(outfile)
	Muxer.PaddingToMakeCounterCont = true
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, err.Error())
		return
	}
	data, err := Storage.StreamHLSTS(uuid, stringToInt(c.Param("seq")))
	if err != nil {
		c.IndentedJSON(500, err.Error())
		return
	}
	if len(data) == 0 {
		c.IndentedJSON(500, "No Segment Found")
		return
	}
	for _, v := range data {
		v.CompositionTime = 1
		err = Muxer.WritePacket(*v)
		if err != nil {
			c.IndentedJSON(500, err.Error())
			return
		}
	}
	err = Muxer.WriteTrailer()
	if err != nil {
		c.IndentedJSON(500, err.Error())
		return
	}
	_, err = c.Writer.Write(outfile.Bytes())
	if err != nil {
		c.IndentedJSON(400, err.Error())
		return
	}

}
