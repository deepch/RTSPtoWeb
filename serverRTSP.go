package main

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	OPTIONS  = "OPTIONS"
	DESCRIBE = "DESCRIBE"
	SETUP    = "SETUP"
	PLAY     = "PLAY"
	TEARDOWN = "TEARDOWN"
)

func RTSPServer() {
	log.WithFields(logrus.Fields{
		"module": "rtsp_server",
		"func":   "RTSPServer",
		"call":   "Start",
	}).Infoln("Server RTSP start")
	l, err := net.Listen("tcp", Storage.ServerRTSPPort())
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "rtsp_server",
			"func":   "RTSPServer",
			"call":   "Listen",
		}).Errorln(err)
		os.Exit(1)
	}
	defer func() {
		err := l.Close()
		if err != nil {
			log.WithFields(logrus.Fields{
				"module": "rtsp_server",
				"func":   "RTSPServer",
				"call":   "Close",
			}).Errorln(err)
		}
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.WithFields(logrus.Fields{
				"module": "rtsp_server",
				"func":   "RTSPServer",
				"call":   "Accept",
			}).Errorln(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	buf := make([]byte, 4096)
	uuid, cid, channel, in, cSEQ := "", "", 0, 0, 0
	var ch chan *[]byte
	defer func() {
		err := conn.Close()
		if err != nil {
			log.WithFields(logrus.Fields{
				"module":  "rtsp_server",
				"stream":  uuid,
				"channel": channel,
				"func":    "handleRequest",
				"call":    "Close",
			}).Errorln(err.Error())
		}
		Storage.ClientDelete(uuid, cid, channel)
		log.WithFields(logrus.Fields{
			"module":  "rtsp_server",
			"stream":  uuid,
			"channel": channel,
			"func":    "handleRequest",
			"call":    "ClientDelete",
		}).Infoln("Client offline")
	}()
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.WithFields(logrus.Fields{
			"module":  "rtsp_server",
			"stream":  uuid,
			"channel": channel,
			"func":    "handleRequest",
			"call":    "SetDeadline",
		}).Errorln(err.Error())
		return
	}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.WithFields(logrus.Fields{
				"module":  "rtsp_server",
				"stream":  uuid,
				"channel": channel,
				"func":    "handleRequest",
				"call":    "Read",
			}).Errorln(err.Error())
			return
		}
		cSEQ = parsecSEQ(buf[:n])
		stage, err := parseStage(buf[:n])
		if err != nil {
			log.WithFields(logrus.Fields{
				"module":  "rtsp_server",
				"stream":  uuid,
				"channel": channel,
				"func":    "handleRequest",
				"call":    "parseStage",
			}).Errorln(err.Error())
		}
		switch stage {
		case OPTIONS:
			uuid, channel, err = parseStreamChannel(buf[:n])
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "parseStreamChannel",
				}).Errorln(err.Error())
				return
			}
			if !Storage.StreamChannelExist(uuid, channel) {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "StreamChannelExist",
				}).Errorln(ErrorStreamNotFound.Error())
				_, err := conn.Write([]byte("RTSP/1.0 404 Not Found\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\n\r\n"))
				if err != nil {
					log.WithFields(logrus.Fields{
						"module":  "rtsp_server",
						"stream":  uuid,
						"channel": channel,
						"func":    "handleRequest",
						"call":    "Write",
					}).Errorln(err.Error())
					return
				}
				return
			}
			Storage.StreamRun(uuid, channel)
			cid, _, ch, err = Storage.ClientAdd(uuid, channel, RTSP)
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "ClientAdd",
				}).Errorln(err.Error())
				return
			}
			_, err := conn.Write([]byte("RTSP/1.0 200 OK\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\nPublic: DESCRIBE, SETUP, TEARDOWN, PLAY\r\n\r\n"))
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "Write",
				}).Errorln(err.Error())
				return
			}
		case SETUP:
			if !strings.Contains(string(buf[:n]), "interleaved") {
				_, err = conn.Write([]byte("RTSP/1.0 461 Unsupported transport\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\n\r\n"))
				if err != nil {
					log.WithFields(logrus.Fields{
						"module":  "rtsp_server",
						"stream":  uuid,
						"channel": channel,
						"func":    "handleRequest",
						"call":    "Write",
					}).Errorln(err.Error())
					return
				}
				continue
			}
			_, err = conn.Write([]byte("RTSP/1.0 200 OK\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\nUser-Agent: Lavf58.29.100\nSession: aaaaaaa\r\nTransport: RTP/AVP/TCP;unicast;interleaved=" + strconv.Itoa(in) + "-" + strconv.Itoa(in+1) + "\r\n\r\n"))
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "Write",
				}).Errorln(err.Error())
				return
			}
			in = in + 2
		case DESCRIBE:
			sdp, err := Storage.StreamSDP(uuid, channel)
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "StreamSDP",
				}).Errorln(err.Error())
				return
			}
			_, err = conn.Write(append([]byte("RTSP/1.0 200 OK\r\nCSeq: "+strconv.Itoa(cSEQ)+"\r\nUser-Agent: Lavf58.29.100\r\nSession: aaaaaaa\r\nContent-Type: application/sdp\r\nContent-Length: "+strconv.Itoa(len(sdp))+"\r\n\r\n"), sdp...))
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "Write",
				}).Errorln(err.Error())
				return
			}
		case PLAY:
			_, err = conn.Write([]byte("RTSP/1.0 200 OK\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\nUser-Agent: Lavf58.29.100\r\nSession: aaaaaaa\r\n\r\n"))
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "Write",
				}).Errorln(err.Error())
				return
			}
			noVideo := time.NewTimer(10 * time.Second)
			for {
				select {
				case <-noVideo.C:
					return
				case pck := <-ch:
					noVideo.Reset(10 * time.Second)
					err := conn.SetDeadline(time.Now().Add(10 * time.Second))
					if err != nil {
						log.WithFields(logrus.Fields{
							"module":  "rtsp_server",
							"stream":  uuid,
							"channel": channel,
							"func":    "handleRequest",
							"call":    "SetDeadline",
						}).Errorln(err.Error())
						return
					}
					_, err = conn.Write(*pck)
					if err != nil {
						log.WithFields(logrus.Fields{
							"module":  "rtsp_server",
							"stream":  uuid,
							"channel": channel,
							"func":    "handleRequest",
							"call":    "Write",
						}).Errorln(err.Error())
						return
					}
				}
			}
		case TEARDOWN:
			_, err := conn.Write([]byte("RTSP/1.0 200 OK\r\nCSeq: " + strconv.Itoa(cSEQ) + "\r\n\r\n"))
			if err != nil {
				log.WithFields(logrus.Fields{
					"module":  "rtsp_server",
					"stream":  uuid,
					"channel": channel,
					"func":    "handleRequest",
					"call":    "Write",
				}).Errorln(err.Error())
			}
			return
		default:
			log.WithFields(logrus.Fields{
				"module":  "rtsp_server",
				"stream":  uuid,
				"channel": channel,
				"func":    "handleRequest",
				"call":    "Stage",
			}).Errorln("stage bad", stage)
		}
	}
}

//parsecSEQ
func parsecSEQ(buf []byte) int {
	return stringToInt(stringInBetween(string(buf), "CSeq: ", "\r\n"))
}

//parseStage func
func parseStage(buf []byte) (string, error) {
	st := strings.Split(string(buf), " ")
	if len(st) > 0 {
		return st[0], nil
	}
	return "", errors.New("parse stage error " + string(buf))
}

//parseStreamChannel func
func parseStreamChannel(buf []byte) (string, int, error) {
	uri := stringInBetween(string(buf), " ", " ")
	st := strings.Split(uri, "/")
	if len(st) >= 5 {
		return st[3], stringToInt(st[4]), nil
	}
	return "", 0, errors.New("parse stream error " + string(buf))
}
