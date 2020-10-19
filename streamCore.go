package main

import (
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/sirupsen/logrus"
)

//StreamServerRunStreamDo stream run do mux
func StreamServerRunStreamDo(streamID string, channelID int) {
	var status int
	defer func() {
		//TODO fix it no need unlock run if delete stream
		if status != 2 {
			Storage.StreamChannelUnlock(streamID, channelID)
		}
	}()
	for {
		log.WithFields(logrus.Fields{
			"module":  "core",
			"stream":  streamID,
			"channel": channelID,
			"func":    "StreamServerRunStreamDo",
			"call":    "Run",
		}).Infoln("Run stream")

		opt, err := Storage.StreamChannelControl(streamID, channelID)
		if opt.OnDemand && !Storage.ClientHas(streamID, channelID) {
			log.WithFields(logrus.Fields{
				"module":  "core",
				"stream":  streamID,
				"channel": channelID,
				"func":    "StreamServerRunStreamDo",
				"call":    "ClientHas",
			}).Infoln("Stop stream no client")
			return
		}
		if err != nil {
			log.WithFields(logrus.Fields{
				"module":  "core",
				"stream":  streamID,
				"channel": channelID,
				"func":    "StreamServerRunStreamDo",
				"call":    "Restart",
			}).Infoln("Restart stream", err)
		}
		status, err = StreamServerRunStream(streamID, channelID, opt)
		if status > 0 {
			log.WithFields(logrus.Fields{
				"module":  "core",
				"stream":  streamID,
				"channel": channelID,
				"func":    "StreamServerRunStreamDo",
				"call":    "StreamServerRunStream",
			}).Infoln("Stream exit by signal or not client")
			return
		}
		if err != nil {
			log.WithFields(logrus.Fields{
				"module":  "core",
				"stream":  streamID,
				"channel": channelID,
				"func":    "StreamServerRunStreamDo",
				"call":    "Restart",
			}).Errorln("Stream error restart stream", err)
		}
		time.Sleep(2 * time.Second)

	}
}

//StreamServerRunStream core stream
func StreamServerRunStream(streamID string, channelID int, opt *ChannelST) (int, error) {
	keyTest := time.NewTimer(20 * time.Second)
	checkClients := time.NewTimer(20 * time.Second)
	var preKeyTS = time.Duration(0)
	var Seq []*av.Packet
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: opt.URL, DisableAudio: true, DialTimeout: 3 * time.Second, ReadWriteTimeout: time.Second * 5 * time.Second, Debug: opt.Debug})
	if err != nil {
		return 0, err
	}
	Storage.StreamChannelStatus(streamID, channelID, ONLINE)
	defer func() {
		RTSPClient.Close()
		Storage.StreamChannelStatus(streamID, channelID, OFFLINE)
		Storage.StreamHLSFlush(streamID, channelID)
	}()
	if len(RTSPClient.CodecData) > 0 {
		Storage.StreamChannelCodecsUpdate(streamID, channelID, RTSPClient.CodecData, RTSPClient.SDPRaw)
	}
	log.WithFields(logrus.Fields{
		"module":  "core",
		"stream":  streamID,
		"channel": channelID,
		"func":    "StreamServerRunStream",
		"call":    "Start",
	}).Infoln("Success connection RTSP")
	for {
		select {
		//Check stream have clients
		case <-checkClients.C:
			if opt.OnDemand && !Storage.ClientHas(streamID, channelID) {
				return 1, ErrorStreamNoClients
			}
			checkClients.Reset(20 * time.Second)
		//Check stream send key
		case <-keyTest.C:
			return 0, ErrorStreamNoVideo
		//Read core signals
		case signals := <-opt.signals:
			switch signals {
			case SignalStreamStop:
				return 2, ErrorStreamStopCoreSignal
			case SignalStreamRestart:
				return 0, ErrorStreamRestart
			case SignalStreamClient:
				return 1, ErrorStreamNoClients
			}
		//Read rtsp signals
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				Storage.StreamChannelCodecsUpdate(streamID, channelID, RTSPClient.CodecData, RTSPClient.SDPRaw)
			case rtspv2.SignalStreamRTPStop:
				return 0, ErrorStreamStopRTSPSignal
			}
		case packetRTP := <-RTSPClient.OutgoingProxy:
			keyTest.Reset(20 * time.Second)
			Storage.StreamChannelCastProxy(streamID, channelID, packetRTP)
		case packetAV := <-RTSPClient.OutgoingPacket:
			if packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
				if preKeyTS > 0 {
					Storage.StreamHLSAdd(streamID, channelID, Seq, packetAV.Time-preKeyTS)
					Seq = []*av.Packet{}
				}
				preKeyTS = packetAV.Time
			}
			Seq = append(Seq, packetAV)
			Storage.StreamChannelCast(streamID, channelID, packetAV)
		}
	}
}
