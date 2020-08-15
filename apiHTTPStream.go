package main

import (
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreams function return stream list
func HTTPAPIServerStreams(c *gin.Context) {
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.List()})
}

//HTTPAPIServerStreamAdd function add new stream
func HTTPAPIServerStreamAdd(c *gin.Context) {
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	err = Storage.StreamAdd(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamEdit function edit stream
func HTTPAPIServerStreamEdit(c *gin.Context) {
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	err = Storage.StreamEdit(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamDelete function delete stream
func HTTPAPIServerStreamDelete(c *gin.Context) {
	err := Storage.StreamDelete(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamDelete function reload stream
func HTTPAPIServerStreamReload(c *gin.Context) {
	err := Storage.StreamReload(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamInfo function return stream info struct
func HTTPAPIServerStreamInfo(c *gin.Context) {
	info, err := Storage.StreamInfo(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: info})
}

//HTTPAPIServerStreamCodec function return codec info struct
func HTTPAPIServerStreamCodec(c *gin.Context) {
	if !Storage.StreamExist(c.Param("uuid")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		return
	}
	codecs, err := Storage.StreamCodecs(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		loggingPrintln(c.Param("uuid"), Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: codecs})
}
