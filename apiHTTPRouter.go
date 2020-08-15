package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	/*
		Stream Control elements
	*/
	privat.GET("/streams", HTTPAPIServerStreams)
	privat.POST("/stream/:uuid/add", HTTPAPIServerStreamAdd)
	privat.POST("/stream/:uuid/edit", HTTPAPIServerStreamEdit)
	privat.GET("/stream/:uuid/delete", HTTPAPIServerStreamDelete)
	privat.GET("/stream/:uuid/reload", HTTPAPIServerStreamReload)
	privat.GET("/stream/:uuid/info", HTTPAPIServerStreamInfo)
	privat.GET("/stream/:uuid/codec", HTTPAPIServerStreamCodec)
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
	/*
		Static HTML Files Demo Mode
	*/
	if Storage.ServerHTTPDemo() {
		public.StaticFS("/static", http.Dir("web/static"))
	}
	err := public.Run(Storage.ServerHTTPPort())
	if err != nil {
		loggingPrintln(Message{Status: 0, Payload: err.Error()})
		os.Exit(1)
	}
}

//HTTPAPIServerIndex index file
func HTTPAPIServerIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
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
