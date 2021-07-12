package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreamChannelCodec function return codec info struct
func HTTPAPIServerStreamChannelCodec(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_stream",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamChannelCodec",
	})

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelCodec",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: codecs})
}

//HTTPAPIServerStreamChannelInfo function return stream info struct
func HTTPAPIServerStreamChannelInfo(c *gin.Context) {
	info, err := Storage.StreamChannelInfo(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module":  "http_stream",
			"stream":  c.Param("uuid"),
			"channel": c.Param("channel"),
			"func":    "HTTPAPIServerStreamChannelInfo",
			"call":    "StreamChannelInfo",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: info})
}

//HTTPAPIServerStreamChannelReload function reload stream
func HTTPAPIServerStreamChannelReload(c *gin.Context) {
	err := Storage.StreamChannelReload(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module":  "http_stream",
			"stream":  c.Param("uuid"),
			"channel": c.Param("channel"),
			"func":    "HTTPAPIServerStreamChannelReload",
			"call":    "StreamChannelReload",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamChannelEdit function edit stream
func HTTPAPIServerStreamChannelEdit(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_stream",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamChannelEdit",
	})

	var payload ChannelST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "BindJSON",
		}).Errorln(err.Error())
		return
	}
	err = Storage.StreamChannelEdit(c.Param("uuid"), c.Param("channel"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelEdit",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamChannelDelete function delete stream
func HTTPAPIServerStreamChannelDelete(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_stream",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamChannelDelete",
	})

	err := Storage.StreamChannelDelete(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelDelete",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamChannelAdd function add new stream
func HTTPAPIServerStreamChannelAdd(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_stream",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamChannelAdd",
	})

	var payload ChannelST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "BindJSON",
		}).Errorln(err.Error())
		return
	}
	err = Storage.StreamChannelAdd(c.Param("uuid"), c.Param("channel"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelAdd",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}
