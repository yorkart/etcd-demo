package etcd

import (
	"sync"
)

var conn *ETCDConn
var rwLock sync.RWMutex

func Set(c *ETCDConn) {
	rwLock.Lock()
	conn = c
	rwLock.Unlock()
}

func Get() *ETCDConn {
	rwLock.RLock()
	c := conn
	rwLock.RUnlock()

	return c
}

func CreateConn() error {
	caCertPath := "/data/tmp/cert/test/ca.crt"
	clientCertPath := "/data/tmp/cert/test/client.crt"
	clientKeyPath :=  "/data/tmp/cert/test/client.key"

	endpoints := []string{"https://10.14.91.12:2379","https://10.14.91.13:2379",}

	c, err := NewETCDConnWithCertPath(caCertPath, clientCertPath, clientKeyPath, endpoints)
	if err != nil {
		return err
	}

	Set(c)

	return nil
}
