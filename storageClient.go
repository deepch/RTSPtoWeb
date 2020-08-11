package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//ready
//ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(uuid string) (string, chan *av.Packet, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	tmp, ok := obj.Streams[uuid]
	if !ok {
		return "", nil, ErrorNotFound
	}
	//Generate UUID client
	cid, err := generateUUID()
	if err != nil {
		return "", nil, err
	}
	ch := make(chan *av.Packet, 2000)
	tmp.clients[cid] = ClientST{outgoingPacket: ch}
	tmp.ack = time.Now()
	obj.Streams[uuid] = tmp
	return cid, ch, nil

}

//ready
//ClientDelete Delete Client
func (obj *StorageST) ClientDelete(uuid string, cid string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; ok {
		delete(obj.Streams[uuid].clients, cid)
	}
}

func (obj *StorageST) ClientHas(uuid string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	tmp, ok := obj.Streams[uuid]
	if !ok {
		return false
	}
	if time.Now().Sub(tmp.ack).Seconds() > 30 {
		return false
	}
	return true
}
