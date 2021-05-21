package main

import (
	"sort"
	"strconv"
	"time"

	"github.com/deepch/vdk/av"
)

//StreamHLSAdd add hls seq to buffer
func (obj *StorageST) StreamHLSAdd(uuid string, channelID string, val []*av.Packet, dur time.Duration) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsSegmentNumber++
			channelTmp.hlsSegmentBuffer[channelTmp.hlsSegmentNumber] = SegmentOld{data: val, dur: dur}
			if len(channelTmp.hlsSegmentBuffer) >= 6 {
				delete(channelTmp.hlsSegmentBuffer, channelTmp.hlsSegmentNumber-6-1)
			}
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}

//StreamHLSm3u8 get hls m3u8 list
func (obj *StorageST) StreamHLSm3u8(uuid string, channelID string) (string, int, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			var out string
			//TODO fix  it
			out += "#EXTM3U\r\n#EXT-X-TARGETDURATION:4\r\n#EXT-X-VERSION:4\r\n#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(channelTmp.hlsSegmentNumber) + "\r\n"
			var keys []int
			for k := range channelTmp.hlsSegmentBuffer {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			var count int
			for _, i := range keys {
				count++
				out += "#EXTINF:" + strconv.FormatFloat(channelTmp.hlsSegmentBuffer[i].dur.Seconds(), 'f', 1, 64) + ",\r\nsegment/" + strconv.Itoa(i) + "/file.ts\r\n"

			}
			return out, count, nil
		}
	}
	return "", 0, ErrorStreamNotFound
}

//StreamHLSTS send hls segment buffer to clients
func (obj *StorageST) StreamHLSTS(uuid string, channelID string, seq int) ([]*av.Packet, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if tmp, ok := channelTmp.hlsSegmentBuffer[seq]; ok {
				return tmp.data, nil
			}
		}
	}
	return nil, ErrorStreamNotFound
}

//StreamHLSFlush delete hls cache
func (obj *StorageST) StreamHLSFlush(uuid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsSegmentBuffer = make(map[int]SegmentOld)
			channelTmp.hlsSegmentNumber = 0
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}
