package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//Fragment struct
type Fragment struct {
	Independent bool          //Fragment have i-frame (key frame)
	Finish      bool          //Fragment Ready
	Duration    time.Duration //Fragment Duration
	Packets     []*av.Packet  //Packet Slice
}

//NewFragment open new fragment
func (element *Segment) NewFragment() *Fragment {
	res := &Fragment{}
	element.Fragment[element.CurrentFragmentID] = res
	return res
}

//GetDuration return fragment dur
func (element *Fragment) GetDuration() time.Duration {
	return element.Duration
}

//WritePacket to fragment func
func (element *Fragment) WritePacket(packet *av.Packet) {
	//increase fragment dur
	element.Duration += packet.Duration
	//Independent if have key
	if packet.IsKeyFrame {
		element.Independent = true
	}
	//append packet to slice of packet
	element.Packets = append(element.Packets, packet)
}

//Close fragment block func
func (element *Fragment) Close() {
	//TODO add callback func
	//finalize fragment
	element.Finish = true
}
