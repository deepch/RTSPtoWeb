package main

/*
 Server Sections
*/

func (obj *StorageST) ServerHTTPDebug() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDebug
}

func (obj *StorageST) ServerHTTPDemo() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDemo
}

func (obj *StorageST) ServerHTTPLogin() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPLogin
}

func (obj *StorageST) ServerHTTPPassword() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPassword
}

func (obj *StorageST) ServerHTTPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPort
}
