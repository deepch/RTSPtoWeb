package main

import (
	"github.com/deepch/vdk/format/mp4f"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreamHLSLLInit send client ts segment
func HTTPAPIServerStreamHLSLLInit(c *gin.Context) {
	safeContext := c.Copy()
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hlsll",
		"stream":  safeContext.Param("uuid"),
		"channel": safeContext.Param("channel"),
		"func":    "HTTPAPIServerStreamHLSLLInit",
	})

	if !Storage.StreamChannelExist(safeContext.Param("uuid"), safeContext.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("HLS", safeContext.Param("uuid"), safeContext.Param("channel"), safeContext.Query("token"), safeContext.ClientIP()) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamUnauthorized.Error())
		return
	}

	c.Header("Content-Type", "application/x-mpegURL")
	Storage.StreamChannelRun(safeContext.Param("uuid"), safeContext.Param("channel"))
	codecs, err := Storage.StreamChannelCodecs(safeContext.Param("uuid"), safeContext.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelCodecs",
		}).Errorln(err.Error())
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	c.Header("Content-Type", "video/mp4")
	_, buf := Muxer.GetInit(codecs)
	_, err = c.Writer.Write(buf)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}
}

//HTTPAPIServerStreamHLSLLM3U8 send client m3u8 play list
func HTTPAPIServerStreamHLSLLM3U8(c *gin.Context) {
	safeContext := c.Copy()
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hlsll",
		"stream":  safeContext.Param("uuid"),
		"channel": safeContext.Param("channel"),
		"func":    "HTTPAPIServerStreamHLSLLM3U8",
	})

	if !Storage.StreamChannelExist(safeContext.Param("uuid"), safeContext.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	c.Header("Content-Type", "application/x-mpegURL")
	Storage.StreamChannelRun(safeContext.Param("uuid"), safeContext.Param("channel"))
	index, err := Storage.HLSMuxerM3U8(safeContext.Param("uuid"), safeContext.Param("channel"), stringToInt(safeContext.DefaultQuery("_HLS_msn", "-1")), stringToInt(safeContext.DefaultQuery("_HLS_part", "-1")))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "HLSMuxerM3U8",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	_, err = c.Writer.Write([]byte(index))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
}

//HTTPAPIServerStreamHLSLLM4Segment send client ts segment
func HTTPAPIServerStreamHLSLLM4Segment(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hlsll",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamHLSLLM4Segment",
	})

	c.Header("Content-Type", "video/mp4")
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
			"call": "StreamChannelCodecs",
		}).Errorln(err.Error())
		return
	}
	if codecs == nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln("Codec Null")
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	seqData, err := Storage.HLSMuxerSegment(c.Param("uuid"), c.Param("channel"), stringToInt(c.Param("segment")))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "HLSMuxerSegment",
		}).Errorln(err.Error())
		return
	}
	for _, v := range seqData {
		err = Muxer.WritePacket4(*v)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "WritePacket4",
			}).Errorln(err.Error())
			return
		}
	}
	buf := Muxer.Finalize()
	_, err = c.Writer.Write(buf)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}
}

//HTTPAPIServerStreamHLSLLM4Fragment send client ts segment
func HTTPAPIServerStreamHLSLLM4Fragment(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hlsll",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamHLSLLM4Fragment",
	})

	c.Header("Content-Type", "video/mp4")
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
			"call": "StreamChannelCodecs",
		}).Errorln(err.Error())
		return
	}
	if codecs == nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln("Codec Null")
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	seqData, err := Storage.HLSMuxerFragment(c.Param("uuid"), c.Param("channel"), stringToInt(c.Param("segment")), stringToInt(c.Param("fragment")))
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "HLSMuxerFragment",
		}).Errorln(err.Error())
		return
	}
	for _, v := range seqData {
		err = Muxer.WritePacket4(*v)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"call": "WritePacket4",
			}).Errorln(err.Error())
			return
		}
	}
	buf := Muxer.Finalize()
	_, err = c.Writer.Write(buf)
	if err != nil {
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}
}
