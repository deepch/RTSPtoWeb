package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//Segment struct
type Segment struct {
	FPS               int               //Current fps
	CurrentFragment   *Fragment         //CurrentFragment link
	CurrentFragmentID int               //CurrentFragment ID
	Finish            bool              //Segment Ready
	Duration          time.Duration     //Segment Duration
	Time              time.Time         //Realtime EXT-X-PROGRAM-DATE-TIME
	Fragment          map[int]*Fragment //Fragment map
}

//NewSegment func
func (element *MuxerHLS) NewSegment() *Segment {
	res := &Segment{
		Fragment:          make(map[int]*Fragment),
		CurrentFragmentID: -1, //Default fragment -1
	}
	//Increase MSN
	element.MSN++
	element.Segments[element.MSN] = res
	return res
}

//GetDuration func
func (element *Segment) GetDuration() time.Duration {
	return element.Duration
}

//SetFPS func
func (element *Segment) SetFPS(fps int) {
	element.FPS = fps
}

//WritePacket func
func (element *Segment) WritePacket(packet *av.Packet) {
	if element.CurrentFragment == nil || element.CurrentFragment.GetDuration().Milliseconds() >= element.FragmentMS(element.FPS) {
		if element.CurrentFragment != nil {
			element.CurrentFragment.Close()
		}
		element.CurrentFragmentID++
		element.CurrentFragment = element.NewFragment()
	}
	element.Duration += packet.Duration
	element.CurrentFragment.WritePacket(packet)
}

//GetFragmentID func
func (element *Segment) GetFragmentID() int {
	return element.CurrentFragmentID
}

//Close segment func
func (element *Segment) Close() {
	element.Finish = true
	if element.CurrentFragment != nil {
		element.CurrentFragment.Close()
	}
}

//FragmentMS func
func (element *Segment) FragmentMS(fps int) int64 {
	for i := 6; i >= 1; i-- {
		if fps%i == 0 {
			return int64(float64(1000) / float64(fps) * float64(i))
		}
	}
	return 100
}
