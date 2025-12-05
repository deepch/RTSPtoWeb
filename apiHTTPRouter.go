package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Message resp struct
type Message struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

// HTTPAPIServer start http server routes
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
	// Login routes
	public.POST("/login", HTTPAPILoginSubmit)
	public.GET("/logout", HTTPAPILogout)

	//Add private login password protect methods
	privat := public.Group("/")
	privat.Use(AuthMiddleware())

	admin := privat.Group("/")
	admin.Use(AdminOnly())

	/*
		Serve React App
	*/
	public.Static("/assets", "./frontend/dist/assets")
	public.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	public.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	/*
		Stream Control elements
	*/

	privat.GET("/streams", HTTPAPIServerStreams)
	admin.POST("/stream/:uuid/add", HTTPAPIServerStreamAdd)
	admin.POST("/stream/:uuid/edit", HTTPAPIServerStreamEdit)
	admin.GET("/stream/:uuid/delete", HTTPAPIServerStreamDelete)
	admin.GET("/stream/:uuid/reload", HTTPAPIServerStreamReload)
	privat.GET("/stream/:uuid/info", HTTPAPIServerStreamInfo)

	/*
		Streams Multi Control elements
	*/

	admin.POST("/streams/multi/control/add", HTTPAPIServerStreamsMultiControlAdd)
	admin.POST("/streams/multi/control/delete", HTTPAPIServerStreamsMultiControlDelete)

	/*
		Stream Channel elements
	*/

	admin.POST("/stream/:uuid/channel/:channel/add", HTTPAPIServerStreamChannelAdd)
	admin.POST("/stream/:uuid/channel/:channel/edit", HTTPAPIServerStreamChannelEdit)
	admin.GET("/stream/:uuid/channel/:channel/delete", HTTPAPIServerStreamChannelDelete)
	privat.GET("/stream/:uuid/channel/:channel/codec", HTTPAPIServerStreamChannelCodec)
	admin.GET("/stream/:uuid/channel/:channel/reload", HTTPAPIServerStreamChannelReload)
	privat.GET("/stream/:uuid/channel/:channel/info", HTTPAPIServerStreamChannelInfo)

	/*
		Stream video elements
	*/
	//HLS
	public.GET("/stream/:uuid/channel/:channel/hls/live/index.m3u8", HTTPAPIServerStreamHLSM3U8)
	public.GET("/stream/:uuid/channel/:channel/hls/live/segment/:seq/file.ts", HTTPAPIServerStreamHLSTS)
	//HLS remote record
	//public.GET("/stream/:uuid/channel/:channel/hls/rr/:s/:e/index.m3u8", HTTPAPIServerStreamRRM3U8)
	//public.GET("/stream/:uuid/channel/:channel/hls/rr/:s/:e/:seq/file.ts", HTTPAPIServerStreamRRTS)
	//HLS LL
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/index.m3u8", HTTPAPIServerStreamHLSLLM3U8)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/init.mp4", HTTPAPIServerStreamHLSLLInit)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/segment/:segment/:any", HTTPAPIServerStreamHLSLLM4Segment)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/fragment/:segment/:fragment/:any", HTTPAPIServerStreamHLSLLM4Fragment)
	//MSE
	public.GET("/stream/:uuid/channel/:channel/mse", HTTPAPIServerStreamMSE)
	public.POST("/stream/:uuid/channel/:channel/webrtc", HTTPAPIServerStreamWebRTC)
	//Save fragment to mp4
	public.GET("/stream/:uuid/channel/:channel/save/mp4/fragment/:duration", HTTPAPIServerStreamSaveToMP4)
	/*
		HTTPS Mode Cert
		# Key considerations for algorithm "RSA" ≥ 2048-bit
		openssl genrsa -out server.key 2048

		# Key considerations for algorithm "ECDSA" ≥ secp384r1
		# List ECDSA the supported curves (openssl ecparam -list_curves)
		#openssl ecparam -genkey -name secp384r1 -out server.key
		#Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

		openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	*/
	if Storage.ServerHTTPS() {
		if Storage.ServerHTTPSAutoTLSEnable() {
			go func() {
				err := autotls.Run(public, Storage.ServerHTTPSAutoTLSName()+Storage.ServerHTTPSPort())
				if err != nil {
					log.Println("Start HTTPS Server Error", err)
				}
			}()
		} else {
			go func() {
				err := public.RunTLS(Storage.ServerHTTPSPort(), Storage.ServerHTTPSCert(), Storage.ServerHTTPSKey())
				if err != nil {
					log.WithFields(logrus.Fields{
						"module": "http_router",
						"func":   "HTTPSAPIServer",
						"call":   "ServerHTTPSPort",
					}).Fatalln(err.Error())
					os.Exit(1)
				}
			}()
		}
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



type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HTTPAPILoginSubmit login submit
func HTTPAPILoginSubmit(c *gin.Context) {
	var payload LoginPayload
	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Message{Status: 0, Payload: err.Error()})
		return
	}

	// Check legacy config
	if Storage.ServerHTTPLogin() != "" && Storage.ServerHTTPPassword() != "" {
		if payload.Username == Storage.ServerHTTPLogin() && payload.Password == Storage.ServerHTTPPassword() {
			createSession(c, payload.Username, "admin")
			return
		}
	}

	// Check Users map
	if u, ok := Storage.Server.Users[payload.Username]; ok {
		if u.Password == payload.Password {
			createSession(c, payload.Username, u.Role)
			return
		}
	}

	c.IndentedJSON(http.StatusUnauthorized, Message{Status: 0, Payload: "Invalid credentials"})
}

func createSession(c *gin.Context, username, role string) {
	sessionID := uuid.New().String()
	Storage.Server.Sessions[sessionID] = SessionST{
		Username: username,
		Role:     role,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie("RTSP_SESSION", sessionID, 3600*24, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, Message{Status: 1, Payload: "Success"})
}

// HTTPAPILogout logout
func HTTPAPILogout(c *gin.Context) {
	cookie, err := c.Cookie("RTSP_SESSION")
	if err == nil {
		delete(Storage.Server.Sessions, cookie)
	}
	c.SetCookie("RTSP_SESSION", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}







// CrossOrigin Access-Control-Allow-Origin any methods
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

// AuthMiddleware check user and password
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no auth configured, allow everyone as admin (backward compatibility)
		if Storage.ServerHTTPLogin() == "" && Storage.ServerHTTPPassword() == "" && len(Storage.Server.Users) == 0 {
			c.Set("role", "admin")
			c.Next()
			return
		}

		// Check Cookie
		cookie, err := c.Cookie("RTSP_SESSION")
		if err == nil {
			if session, ok := Storage.Server.Sessions[cookie]; ok {
				if session.Expires.After(time.Now()) {
					c.Set("role", session.Role)
					c.Next()
					return
				} else {
					delete(Storage.Server.Sessions, cookie)
				}
			}
		}

		// Check Basic Auth (API support)
		user, pass, hasAuth := c.Request.BasicAuth()
		if hasAuth {
			// Check legacy config
			if Storage.ServerHTTPLogin() != "" && Storage.ServerHTTPPassword() != "" {
				if user == Storage.ServerHTTPLogin() && pass == Storage.ServerHTTPPassword() {
					c.Set("role", "admin")
					c.Next()
					return
				}
			}
			// Check new Users map
			if u, ok := Storage.Server.Users[user]; ok {
				if u.Password == pass {
					c.Set("role", u.Role)
					c.Next()
					return
				}
			}
		}

		// Redirect to login if HTML request
		if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// AdminOnly check if user is admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
