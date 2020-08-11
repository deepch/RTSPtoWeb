package main

import (
	"errors"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
)

//StreamServerRunStreamDo stream run do mux
func StreamServerRunStreamDo(name string) {
	defer Storage.StreamUnlock(name)
	for {
		loggingPrintln("Run Stream", name)
		exit, err := StreamServerRunStream(name)
		if exit {
			loggingPrintln("Stream Exit by Signal or Not Client")
			return
		}
		if err != nil {
			loggingPrintln("Stream Error", err)
		}
		time.Sleep(2 * time.Second)

	}
}
func StreamServerRunStream(name string) (bool, error) {
	keyTest := time.NewTimer(20 * time.Second)
	checkClients := time.NewTimer(20 * time.Second)
	Control, err := Storage.StreamControl(name)
	if err != nil {
		//TODO fix it
		return true, ErrorNotFound
	}
	var preKeyTS = time.Duration(0)
	var Seq []*av.Packet
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: Control.URL, DisableAudio: true, DialTimeout: 3 * time.Second, ReadWriteTimeout: time.Second * 5 * time.Second, Debug: Control.Debug})
	if err != nil {
		return false, errors.New("RTSP Client Error " + err.Error())
	}
	Storage.StreamStatus(name, ONLINE)
	defer func() {
		RTSPClient.Close()
		Storage.StreamStatus(name, OFFLINE)
		Storage.StreamHLSFlush(name)
	}()

	//if codec data recived
	if len(RTSPClient.CodecData) > 0 {
		Storage.StreamCodecsUpdate(name, RTSPClient.CodecData)
	}
	loggingPrintln("Stream", name, "success connection RTSP")
	for {
		select {
		//Read no video timeout
		case <-checkClients.C:
			if Control.OnDemand && !Storage.ClientHas(name) {
				return true, errors.New("not clients on steam")
			}
			checkClients.Reset(20 * time.Second)
		case <-keyTest.C:
			return false, errors.New("Video Stream No Send Key Frame")
		//Read core signals
		case signals := <-Control.signals:
			switch signals {
			case SignalStreamStop:
				return true, errors.New("Core Stop Signal")
			case SignalStreamRestart:
				return false, errors.New("Core Restart Signal")
			case SignalStreamClient:
				loggingPrintln("New Viwer Signal")
			}
		//Read rtsp signals
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				loggingPrintln("Update Code Info")
				Storage.StreamCodecsUpdate(name, RTSPClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return false, errors.New("RTSP Client Restart Signal")
			}
		case <-RTSPClient.OutgoingProxy:
			//Add Raw Proxy Next Version
		case packet := <-RTSPClient.OutgoingPacket:
			if packet.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
				if preKeyTS > 0 {
					Storage.StreamHLSAdd(name, Seq, packet.Time-preKeyTS)
					Seq = []*av.Packet{}
				}
				preKeyTS = packet.Time
			}
			Seq = append(Seq, packet)
			Storage.Cast(name, packet)
		}
	}
}
