package main

import (
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/sirupsen/logrus"
)

//StreamServerRunStreamDo stream run do mux
func StreamServerRunStreamDo(streamID string, channelID int) {
	defer Storage.StreamUnlock(streamID, channelID)
	for {
		log.WithFields(logrus.Fields{
			"module":  "core",
			"stream":  streamID,
			"channel": channelID,
			"func":    "StreamServerRunStreamDo",
			"call":    "Run",
		}).Infoln("Run stream")

		opt, err := Storage.StreamControl(streamID, channelID)
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
		exit, err := StreamServerRunStream(streamID, channelID, opt)
		if exit {
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
func StreamServerRunStream(streamID string, channelID int, opt *ChannelST) (bool, error) {
	keyTest := time.NewTimer(20 * time.Second)
	checkClients := time.NewTimer(20 * time.Second)
	var preKeyTS = time.Duration(0)
	var Seq []*av.Packet
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: opt.URL, DisableAudio: true, DialTimeout: 3 * time.Second, ReadWriteTimeout: time.Second * 5 * time.Second, Debug: opt.Debug})
	if err != nil {
		return false, err
	}
	Storage.StreamStatus(streamID, channelID, ONLINE)
	defer func() {
		RTSPClient.Close()
		Storage.StreamStatus(streamID, channelID, OFFLINE)
		Storage.StreamHLSFlush(streamID, channelID)
	}()
	if len(RTSPClient.CodecData) > 0 {
		Storage.StreamCodecsUpdate(streamID, channelID, RTSPClient.CodecData, RTSPClient.SDPRaw)
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
				return true, ErrorStreamNoClients
			}
			checkClients.Reset(20 * time.Second)
		//Check stream send key
		case <-keyTest.C:
			return false, ErrorStreamNoVideo
		//Read core signals
		case signals := <-opt.signals:
			switch signals {
			case SignalStreamStop:
				return true, ErrorStreamStopCoreSignal
			case SignalStreamRestart:
				return false, ErrorStreamRestart
			case SignalStreamClient:
				return true, ErrorStreamNoClients
			}
		//Read rtsp signals
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				Storage.StreamCodecsUpdate(streamID, channelID, RTSPClient.CodecData, RTSPClient.SDPRaw)
			case rtspv2.SignalStreamRTPStop:
				return false, ErrorStreamStopRTSPSignal
			}
		case packetRTP := <-RTSPClient.OutgoingProxy:
			keyTest.Reset(20 * time.Second)
			Storage.CastProxy(streamID, channelID, packetRTP)
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
			Storage.Cast(streamID, channelID, packetAV)
		}
	}
}
