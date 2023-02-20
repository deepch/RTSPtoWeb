package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreams function return stream list
func HTTPAPIServerStreams(c *gin.Context) {
	list, err := Storage.MarshalledStreamsList()
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: list})
}

//HTTPAPIServerStreamsMultiControlAdd function add new stream's
func HTTPAPIServerStreamsMultiControlAdd(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module": "http_stream",
		"func":   "HTTPAPIServerStreamsMultiControlAdd",
	})

	var payload StorageST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "BindJSON",
		}).Errorln(err.Error())
		return
	}
	if payload.Streams == nil || len(payload.Streams) < 1 {
		c.IndentedJSON(400, Message{Status: 0, Payload: ErrorStreamsLen0.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "len(payload)",
		}).Errorln(ErrorStreamsLen0.Error())
		return
	}
	var resp = make(map[string]Message)
	var FoundError bool
	for k, v := range payload.Streams {
		err = Storage.StreamAdd(k, v)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"stream": k,
				"call":   "StreamAdd",
			}).Errorln(err.Error())
			resp[k] = Message{Status: 0, Payload: err.Error()}
			FoundError = true
		} else {
			resp[k] = Message{Status: 1, Payload: Success}
		}
	}
	if FoundError {
		c.IndentedJSON(200, Message{Status: 0, Payload: resp})
	} else {
		c.IndentedJSON(200, Message{Status: 1, Payload: resp})
	}
}

//HTTPAPIServerStreamsMultiControlDelete function delete stream's
func HTTPAPIServerStreamsMultiControlDelete(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module": "http_stream",
		"func":   "HTTPAPIServerStreamsMultiControlDelete",
	})

	var payload []string
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "BindJSON",
		}).Errorln(err.Error())
		return
	}
	if len(payload) < 1 {
		c.IndentedJSON(400, Message{Status: 0, Payload: ErrorStreamsLen0.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "len(payload)",
		}).Errorln(ErrorStreamsLen0.Error())
		return
	}
	var resp = make(map[string]Message)
	var FoundError bool
	for _, key := range payload {
		err := Storage.StreamDelete(key)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"stream": key,
				"call":   "StreamDelete",
			}).Errorln(err.Error())
			resp[key] = Message{Status: 0, Payload: err.Error()}
			FoundError = true
		} else {
			resp[key] = Message{Status: 1, Payload: Success}
		}
	}
	if FoundError {
		c.IndentedJSON(200, Message{Status: 0, Payload: resp})
	} else {
		c.IndentedJSON(200, Message{Status: 1, Payload: resp})
	}
}

//HTTPAPIServerStreamAdd function add new stream
func HTTPAPIServerStreamAdd(c *gin.Context) {
	var payload StreamST
	err := c.BindJSON(&payload)
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamAdd",
			"call":   "BindJSON",
		}).Errorln(err.Error())
		return
	}
	err = Storage.StreamAdd(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamAdd",
			"call":   "StreamAdd",
		}).Errorln(err.Error())
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
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamEdit",
			"call":   "BindJSON",
		}).Errorln(err.Error())
		return
	}
	err = Storage.StreamEdit(c.Param("uuid"), payload)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamEdit",
			"call":   "StreamEdit",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamDelete function delete stream
func HTTPAPIServerStreamDelete(c *gin.Context) {
	err := Storage.StreamDelete(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamDelete",
			"call":   "StreamDelete",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamDelete function reload stream
func HTTPAPIServerStreamReload(c *gin.Context) {
	err := Storage.StreamReload(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamReload",
			"call":   "StreamReload",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: Success})
}

//HTTPAPIServerStreamInfo function return stream info struct
func HTTPAPIServerStreamInfo(c *gin.Context) {
	info, err := Storage.StreamInfo(c.Param("uuid"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		log.WithFields(logrus.Fields{
			"module": "http_stream",
			"stream": c.Param("uuid"),
			"func":   "HTTPAPIServerStreamInfo",
			"call":   "StreamInfo",
		}).Errorln(err.Error())
		return
	}
	c.IndentedJSON(200, Message{Status: 1, Payload: info})
}
