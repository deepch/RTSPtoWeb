package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/deepch/vdk/format/ts"
	"github.com/deepch/vdk/format/webrtc"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func HTTPAPIServer() {
	//Set HTTP API mode
	if !Storage.ServerHTTPDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	public := gin.Default()
	public.Use(CrossOrigin())
	//Add private login password protect methods
	privat := public.Group("/", gin.BasicAuth(gin.Accounts{Storage.ServerHTTPLogin(): Storage.ServerHTTPPassword()}))
	public.LoadHTMLGlob("web/templates/*")
	/*
		Html template
	*/
	public.GET("/", HTTPAPIServerIndex)
	/*
		Stream Control elements
	*/
	privat.GET("/streams", HTTPAPIServerStreams)
	privat.POST("/stream/:uuid/add", HTTPAPIServerStreamAdd)
	privat.POST("/stream/:uuid/edit", HTTPAPIServerStreamEdit)
	privat.GET("/stream/:uuid/delete", HTTPAPIServerStreamDelete)
	privat.GET("/stream/:uuid/reload", HTTPAPIServerStreamReload)
	privat.GET("/stream/:uuid/info", HTTPAPIServerStreamInfo)
	/*
		Stream video elements
	*/
	public.GET("/stream/:uuid/hls/live/index.m3u8", HTTPAPIServerStreamHLSM3U8)
	public.GET("/stream/:uuid/hls/live/segment/:seq/file.ts", HTTPAPIServerStreamHLSTS)
	public.GET("/stream/:uuid/mse", func(c *gin.Context) {
		handler := websocket.Handler(HTTPAPIServerStreamMSE)
		handler.ServeHTTP(c.Writer, c.Request)
	})
	public.POST("/stream/:uuid/webrtc", HTTPAPIServerStreamWebRTC)
	//TODO Fix It
	public.GET("/codec/:uuid", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		if Storage.StreamExist(c.Param("uuid")) {
			codecs, _ := Storage.StreamCodecs(c.Param("uuid"))
			if codecs == nil {
				return
			}
			b, err := json.Marshal(codecs)
			log.Println(string(b), err)
			if err == nil {
				_, err = c.Writer.Write(b)
				if err == nil {
					log.Println("Write Codec Info error", err)
					return
				}
			}
		}
	})
	/*
		Static HTML Files Demo Mode
	*/
	if Storage.ServerHTTPDemo() {
		public.StaticFS("/static", http.Dir("web/static"))
	}
	err := public.Run(Storage.ServerHTTPPort())
	if err != nil {
		log.Fatalln(err)
	}
}
func HTTPAPIServerIndex(c *gin.Context) {
	//fi, all := Storage.List()
	//sort.Strings(all)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		//	"uuid":    fi,
		//	"uuidMap": all,
		"version": time.Now().String(),
	})

}

func HTTPAPIServerStreams(c *gin.Context) {
	c.IndentedJSON(200, Storage.List())
}

func HTTPAPIServerStreamAdd(c *gin.Context) {
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	err = Storage.StreamAdd(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamEdit(c *gin.Context) {
	if !Storage.StreamExist(c.Param("uuid")) {
		return
	}
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	err = Storage.StreamEdit(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamDelete(c *gin.Context) {
	err := Storage.StreamDelete(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}
func HTTPAPIServerStreamReload(c *gin.Context) {
	err := Storage.StreamReload(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, "ok")
}

func HTTPAPIServerStreamInfo(c *gin.Context) {
	info, err := Storage.StreamInfo(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, err)
		return
	}
	c.IndentedJSON(200, info)
}

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
	data, err := Storage.StreamHLSTS(uuid, StringToInt(c.Param("seq")))
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

//ready
//CrossOrigin Access-Control-Allow-Origin any methods
func CrossOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
