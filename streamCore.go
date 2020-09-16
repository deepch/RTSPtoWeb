package main

import (
	"log"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
)

//StreamServerRunStreamDo stream run do mux
func StreamServerRunStreamDo(streamID string, channelID int) {
	defer Storage.StreamUnlock(streamID, channelID)
	for {
		loggingPrintln("Run Stream", streamID, "Channel", channelID)
		opt, err := Storage.StreamControl(streamID, channelID)
		if opt.OnDemand && !Storage.ClientHas(streamID, channelID) {
			loggingPrintln("Stream Exit Not Client")
			return
		}
		if err != nil {
			loggingPrintln("Stream Error", err, "Restart Stream")
		}
		exit, err := StreamServerRunStream(streamID, channelID, opt)
		if exit {
			loggingPrintln("Stream Exit by Signal or Not Client")
			return
		}
		if err != nil {
			loggingPrintln("Stream Error", err, "Restart Stream")
		}
		time.Sleep(2 * time.Second)

	}
}

//StreamServerRunStream core stream
func StreamServerRunStream(streamID string, channelID int, opt *ChannelST) (bool, error) {
	log.Println(streamID, "run cid", channelID)
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
		Storage.StreamCodecsUpdate(streamID, channelID, RTSPClient.CodecData)
	}
	loggingPrintln("Stream", streamID, "Channel", channelID, "Success connection RTSP")
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
				Storage.StreamCodecsUpdate(streamID, channelID, RTSPClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return false, ErrorStreamStopRTSPSignal
			}
		case <-RTSPClient.OutgoingProxy:
			//TODO Add Raw Proxy Next Version
		case packet := <-RTSPClient.OutgoingPacket:
			if packet.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
				if preKeyTS > 0 {
					Storage.StreamHLSAdd(streamID, channelID, Seq, packet.Time-preKeyTS)
					Seq = []*av.Packet{}
				}
				preKeyTS = packet.Time
			}
			Seq = append(Seq, packet)
			Storage.Cast(streamID, channelID, packet)
		}
	}
}
