package main

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	//Default www static file dir
	DefaultHTTPDir = "web"
)

//ServerHTTPDir
func (obj *StorageST) ServerHTTPDir() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if filepath.Clean(obj.Server.HTTPDir) == "." {
		return DefaultHTTPDir
	}
	return filepath.Clean(obj.Server.HTTPDir)
}

//ServerHTTPDebug read debug options
func (obj *StorageST) ServerHTTPDebug() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDebug
}

//ServerLogLevel read debug options
func (obj *StorageST) ServerLogLevel() logrus.Level {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.LogLevel
}

//ServerHTTPDemo read demo options
func (obj *StorageST) ServerHTTPDemo() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDemo
}

//ServerHTTPLogin read Login options
func (obj *StorageST) ServerHTTPLogin() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPLogin
}

//ServerHTTPPassword read Password options
func (obj *StorageST) ServerHTTPPassword() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPassword
}

//ServerHTTPPort read HTTP Port options
func (obj *StorageST) ServerHTTPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPort
}

//ServerRTSPPort read HTTP Port options
func (obj *StorageST) ServerRTSPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.RTSPPort
}

//ServerHTTPSPort read HTTPS Port options
func (obj *StorageST) ServerHTTPS() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPS
}

//ServerHTTPSPort read HTTPS Port options
func (obj *StorageST) ServerHTTPSPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSPort
}

//ServerHTTPSCert read HTTPS Cert options
func (obj *StorageST) ServerHTTPSCert() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSCert
}

//ServerHTTPSKey read HTTPS Key options
func (obj *StorageST) ServerHTTPSKey() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSKey
}
