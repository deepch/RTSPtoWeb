package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

//Message resp struct
type Message struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

//HTTPAPIServer start http server routes
func HTTPAPIServer() {
	//Set HTTP API mode
	log.WithFields(logrus.Fields{
		"module": "http_server",
		"func":   "RTSPServer",
		"call":   "Start",
	}).Infoln("Server HTTP start")
	var public *gin.Engine
	if !Storage.ServerHTTPDebug() {
		gin.SetMode(gin.ReleaseMode)
		public = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		public = gin.Default()
	}

	public.Use(CrossOrigin())
	//Add private login password protect methods
	privat := public.Group("/", gin.BasicAuth(gin.Accounts{Storage.ServerHTTPLogin(): Storage.ServerHTTPPassword()}))
	public.LoadHTMLGlob("web/templates/*")
	/*
		Html template
	*/
	public.GET("/", HTTPAPIServerIndex)
	public.GET("/pages/stream/list", HTTPAPIStreamList)
	public.GET("/pages/stream/add", HTTPAPIAddStream)
	public.GET("/pages/stream/edit/:uuid", HTTPAPIEditStream)
	public.GET("/pages/player/hls/:uuid/:channel", HTTPAPIPlayHls)
	public.GET("/pages/player/mse/:uuid/:channel", HTTPAPIPlayMse)
	public.GET("/pages/player/webrtc/:uuid/:channel", HTTPAPIPlayWebrtc)
	public.GET("/pages/multiview", HTTPAPIMultiview)
	public.GET("/pages/fullscreenmulti", HTTPAPIFullScreenMultiview)
	public.POST("/pages/fullscreenmulti", HTTPAPIFullScreenMultiview)
	public.GET("/pages/documentation", HTTPAPIServerDocumentation)
	public.GET("/pages/login", HTTPAPIPageLogin)

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
		Stream Channel elements
	*/

	privat.POST("/stream/:uuid/channel/:channel/add", HTTPAPIServerStreamChannelAdd)
	privat.POST("/stream/:uuid/channel/:channel/edit", HTTPAPIServerStreamChannelEdit)
	privat.POST("/stream/:uuid/channel/:channel/delete", HTTPAPIServerStreamChannelDelete)
	privat.GET("/stream/:uuid/channel/:channel/codec", HTTPAPIServerStreamChannelCodec)
	privat.GET("/stream/:uuid/channel/:channel/reload", HTTPAPIServerStreamChannelReload)
	privat.GET("/stream/:uuid/channel/:channel/info", HTTPAPIServerStreamChannelInfo)

	/*
		Stream video elements
	*/

	public.GET("/stream/:uuid/channel/:channel/hls/live/index.m3u8", HTTPAPIServerStreamHLSM3U8)
	public.GET("/stream/:uuid/channel/:channel/hls/live/segment/:seq/file.ts", HTTPAPIServerStreamHLSTS)
	public.GET("/stream/:uuid/channel/:channel/mse", func(c *gin.Context) {
		handler := websocket.Handler(HTTPAPIServerStreamMSE)
		handler.ServeHTTP(c.Writer, c.Request)
	})
	public.POST("/stream/:uuid/channel/:channel/webrtc", HTTPAPIServerStreamWebRTC)
	/*
		Static HTML Files Demo Mode
	*/
	if Storage.ServerHTTPDemo() {
		public.StaticFS("/static", http.Dir("web/static"))
	}
	err := public.Run(Storage.ServerHTTPPort())
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "http_router",
			"func":   "HTTPAPIServer",
			"call":   "ServerHTTPPort",
		}).Fatalln(err.Error())
		os.Exit(1)
	}
}

//HTTPAPIServerIndex index file
func HTTPAPIServerIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "index",
	})

}

func HTTPAPIServerDocumentation(c *gin.Context) {
	c.HTML(http.StatusOK, "documentation.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "documentation",
	})
}

func HTTPAPIStreamList(c *gin.Context) {
	c.HTML(http.StatusOK, "stream_list.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "stream_list",
	})
}
func HTTPAPIPageLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "login",
	})
}

func HTTPAPIPlayHls(c *gin.Context) {
	c.HTML(http.StatusOK, "play_hls.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_hls",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIPlayMse(c *gin.Context) {
	c.HTML(http.StatusOK, "play_mse.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_mse",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIPlayWebrtc(c *gin.Context) {
	c.HTML(http.StatusOK, "play_webrtc.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_webrtc",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIAddStream(c *gin.Context) {
	c.HTML(http.StatusOK, "add_stream.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "add_stream",
	})
}
func HTTPAPIEditStream(c *gin.Context) {
	c.HTML(http.StatusOK, "edit_stream.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "edit_stream",
		"uuid":    c.Param("uuid"),
	})
}

func HTTPAPIMultiview(c *gin.Context) {
	c.HTML(http.StatusOK, "multiview.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "multiview",
	})
}

func HTTPAPIFullScreenMultiview(c *gin.Context) {
	c.HTML(http.StatusOK, "fullscreenmulti.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "fullscreenmulti",
	})
}
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
