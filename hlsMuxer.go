package main

import (
	"context"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
)

//MuxerHLS struct
type MuxerHLS struct {
	mutex             sync.RWMutex
	UUID              string             //Current UUID
	MSN               int                //Current MSN
	FPS               int                //Current FPS
	MediaSequence     int                //Current MediaSequence
	CurrentFragmentID int                //Current fragment id
	CacheM3U8         string             //Current index cache
	CurrentSegment    *Segment           //Current segment link
	Segments          map[int]*Segment   //Current segments group
	FragmentCtx       context.Context    //chan 1-N
	FragmentCancel    context.CancelFunc //chan 1-N
}

//NewHLSMuxer Segments
func NewHLSMuxer(uuid string) *MuxerHLS {
	ctx, cancel := context.WithCancel(context.Background())
	return &MuxerHLS{
		UUID:           uuid,
		MSN:            -1,
		Segments:       make(map[int]*Segment),
		FragmentCtx:    ctx,
		FragmentCancel: cancel,
	}
}

//SetFPS func
func (element *MuxerHLS) SetFPS(fps int) {
	element.FPS = fps
}

//WritePacket func
func (element *MuxerHLS) WritePacket(packet *av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	//TODO delete packet.IsKeyFrame if need no EXT-X-INDEPENDENT-SEGMENTS

	if !packet.IsKeyFrame && element.CurrentSegment == nil {
		// Wait for the first keyframe before initializing
		return
	}
	if packet.IsKeyFrame && (element.CurrentSegment == nil || element.CurrentSegment.GetDuration().Seconds() >= 4) {
		if element.CurrentSegment != nil {
			element.CurrentSegment.Close()
			if len(element.Segments) > 6 {
				delete(element.Segments, element.MSN-6)
				element.MediaSequence++
			}
		}
		element.CurrentSegment = element.NewSegment()
		element.CurrentSegment.SetFPS(element.FPS)
	}
	element.CurrentSegment.WritePacket(packet)
	CurrentFragmentID := element.CurrentSegment.GetFragmentID()
	if CurrentFragmentID != element.CurrentFragmentID {
		element.UpdateIndexM3u8()
	}
	element.CurrentFragmentID = CurrentFragmentID
}

//UpdateIndexM3u8 func
func (element *MuxerHLS) UpdateIndexM3u8() {
	var header string
	var body string
	var partTarget time.Duration
	var segmentTarget time.Duration
	segmentTarget = time.Second * 2
	for _, segmentKey := range element.SortSegments(element.Segments) {
		for _, fragmentKey := range element.SortFragment(element.Segments[segmentKey].Fragment) {
			if element.Segments[segmentKey].Fragment[fragmentKey].Finish {
				var independent string
				if element.Segments[segmentKey].Fragment[fragmentKey].Independent {
					independent = ",INDEPENDENT=YES"
				}
				body += "#EXT-X-PART:DURATION=" + strconv.FormatFloat(element.Segments[segmentKey].Fragment[fragmentKey].GetDuration().Seconds(), 'f', 5, 64) + "" + independent + ",URI=\"fragment/" + strconv.Itoa(segmentKey) + "/" + strconv.Itoa(fragmentKey) + "/0qrm9ru6." + strconv.Itoa(fragmentKey) + ".m4s\"\n"
				partTarget = element.Segments[segmentKey].Fragment[fragmentKey].Duration
			} else {
				body += "#EXT-X-PRELOAD-HINT:TYPE=PART,URI=\"fragment/" + strconv.Itoa(segmentKey) + "/" + strconv.Itoa(fragmentKey) + "/0qrm9ru6." + strconv.Itoa(fragmentKey) + ".m4s\"\n"
			}
		}
		if element.Segments[segmentKey].Finish {
			segmentTarget = element.Segments[segmentKey].Duration
			body += "#EXT-X-PROGRAM-DATE-TIME:" + element.Segments[segmentKey].Time.Format("2006-01-02T15:04:05.000000Z") + "\n#EXTINF:" + strconv.FormatFloat(element.Segments[segmentKey].Duration.Seconds(), 'f', 5, 64) + ",\n"
			body += "segment/" + strconv.Itoa(segmentKey) + "/" + element.UUID + "." + strconv.Itoa(segmentKey) + ".m4s\n"
		}
	}
	header += "#EXTM3U\n"
	header += "#EXT-X-TARGETDURATION:" + strconv.Itoa(int(math.Round(segmentTarget.Seconds()))) + "\n"
	header += "#EXT-X-VERSION:7\n"
	header += "#EXT-X-INDEPENDENT-SEGMENTS\n"
	header += "#EXT-X-SERVER-CONTROL:CAN-BLOCK-RELOAD=YES,PART-HOLD-BACK=" + strconv.FormatFloat(partTarget.Seconds()*4, 'f', 5, 64) + ",HOLD-BACK=" + strconv.FormatFloat(segmentTarget.Seconds()*4, 'f', 5, 64) + "\n"
	header += "#EXT-X-MAP:URI=\"init.mp4\"\n"
	header += "#EXT-X-PART-INF:PART-TARGET=" + strconv.FormatFloat(partTarget.Seconds(), 'f', 5, 64) + "\n"
	header += "#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(element.MediaSequence) + "\n"
	header += body
	element.CacheM3U8 = header
	element.PlaylistUpdate()
}

//PlaylistUpdate func
func (element *MuxerHLS) PlaylistUpdate() {
	element.FragmentCancel()
	element.FragmentCtx, element.FragmentCancel = context.WithCancel(context.Background())
}

//GetSegment func
func (element *MuxerHLS) GetSegment(segment int) ([]*av.Packet, error) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if segmentTmp, ok := element.Segments[segment]; ok && len(segmentTmp.Fragment) > 0 {
		var res []*av.Packet
		for _, v := range element.SortFragment(segmentTmp.Fragment) {
			res = append(res, segmentTmp.Fragment[v].Packets...)
		}
		return res, nil
	}
	return nil, ErrorStreamNotFound
}

//GetFragment func
func (element *MuxerHLS) GetFragment(segment int, fragment int) ([]*av.Packet, error) {
	element.mutex.Lock()
	if segmentTmp, segmentTmpOK := element.Segments[segment]; segmentTmpOK {
		if fragmentTmp, fragmentTmpOK := segmentTmp.Fragment[fragment]; fragmentTmpOK {
			if fragmentTmp.Finish {
				element.mutex.Unlock()
				return fragmentTmp.Packets, nil
			} else {
				element.mutex.Unlock()
				pck, err := element.WaitFragment(time.Second*1, segment, fragment)
				if err != nil {
					return nil, err
				}
				return pck, err
			}
		}
	}
	element.mutex.Unlock()
	return nil, ErrorStreamNotFound
}

//GetIndexM3u8 func
func (element *MuxerHLS) GetIndexM3u8(needMSN int, needPart int) (string, error) {
	element.mutex.Lock()
	if len(element.CacheM3U8) != 0 && ((needMSN == -1 || needPart == -1) || (needMSN-element.MSN > 1) || (needMSN == element.MSN && needPart < element.CurrentFragmentID)) {
		element.mutex.Unlock()
		return element.CacheM3U8, nil
	} else {
		element.mutex.Unlock()
		index, err := element.WaitIndex(time.Second*3, needMSN, needPart)
		if err != nil {
			return "", err
		}
		return index, err
	}
}

//WaitFragment func
func (element *MuxerHLS) WaitFragment(timeOut time.Duration, segment, fragment int) ([]*av.Packet, error) {
	select {
	case <-time.After(timeOut):
		return nil, ErrorStreamNotFound
	case <-element.FragmentCtx.Done():
		element.mutex.Lock()
		defer element.mutex.Unlock()
		if segmentTmp, segmentTmpOK := element.Segments[segment]; segmentTmpOK {
			if fragmentTmp, fragmentTmpOK := segmentTmp.Fragment[fragment]; fragmentTmpOK {
				if fragmentTmp.Finish {
					return fragmentTmp.Packets, nil
				}
			}
		}
		return nil, ErrorStreamNotFound
	}
}

//WaitIndex func
func (element *MuxerHLS) WaitIndex(timeOut time.Duration, segment, fragment int) (string, error) {
	for {
		select {
		case <-time.After(timeOut):
			return "", ErrorStreamNotFound
		case <-element.FragmentCtx.Done():
			element.mutex.Lock()
			if element.MSN < segment || (element.MSN == segment && element.CurrentFragmentID < fragment) {
				log.Println("wait req", element.MSN, element.CurrentFragmentID, segment, fragment)
				element.mutex.Unlock()
				continue
			}
			element.mutex.Unlock()
			return element.CacheM3U8, nil
		}
	}
}

//SortFragment func
func (element *MuxerHLS) SortFragment(val map[int]*Fragment) []int {
	keys := make([]int, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

//SortSegments fuc
func (element *MuxerHLS) SortSegments(val map[int]*Segment) []int {
	keys := make([]int, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func (element *MuxerHLS) Close() {

}
