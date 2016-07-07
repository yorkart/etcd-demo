package etcd3

import (
	"sync"
)

var (
	endpoints      = []string{"https://10.14.86.144:2389"}
	caCertPath = "/data/tmp/cert/test/ca.crt"
	clientCrtPath = "/data/tmp/cert/test/client.crt"
	clientKeyPath =  "/data/tmp/cert/test/client.key"
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
	c, err := NewETCDConnWithCertPath(caCertPath, clientCrtPath, clientKeyPath, endpoints)
	if err != nil {
		return err
	}

	Set(c)

	return nil
}


//func GetClient(auth *Auth3) (*clientv3.Client, error) {
//	caCrt, err := ioutil.ReadFile(caCertPath)
//	if err != nil {
//		log.Println("ReadFile err:", err)
//		return nil, err
//	}
//	clientCrt, err := ioutil.ReadFile(clientCrtPath)
//	if err != nil {
//		log.Println("ReadFile err:", err)
//		return nil, err
//	}
//	clientKey, err := ioutil.ReadFile(clientKeyPath)
//	if err != nil {
//		log.Println("ReadFile err:", err)
//		return nil, err
//	}
//
//	pool := x509.NewCertPool()
//	pool.AppendCertsFromPEM(caCrt)
//
//	cliCrt, err := tls.X509KeyPair(clientCrt, clientKey)
//	if err != nil {
//		log.Println("Loadx509keypair err:", err)
//		return nil, err
//	}
//
//	t := &tls.Config{
//		RootCAs:      pool,
//		Certificates: []tls.Certificate{cliCrt},
//		//InsecureSkipVerify: true,
//	}
//
//	if auth != nil {
//		cfg := clientv3.Config{
//			Endpoints: endpoints, // []string{"https://10.14.87.171:2379"},
//			TLS: t,
//			// set timeout per request to fail fast when the target endpoint is unavailable
//			DialTimeout: time.Second * 5,
//			Username:                auth.UserName, // "root",
//			Password:                auth.Password, // "123456",
//		}
//		return clientv3.New(cfg)
//	} else {
//		cfg := clientv3.Config{
//			Endpoints: endpoints, // []string{"https://10.14.87.171:2379"},
//			TLS: t,
//			// set timeout per request to fail fast when the target endpoint is unavailable
//			DialTimeout: dialTimeout,
//		}
//		return clientv3.New(cfg)
//	}
//}