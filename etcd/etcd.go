package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/client"
)

type Auth struct {
	UserName string
	Password string
}

type ETCDConn struct {
	transport *http.Transport
	endpoints []string
}

func NewETCDConnWithCertPath(caCertPath string, clientCrtPath string, clientKeyPath string, endpoints []string) (*ETCDConn, error) {
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Println("ReadFile err:", err)
		return nil, err
	}
	certPEMBlock, err := ioutil.ReadFile(clientCrtPath)
	if err != nil {
		log.Println("ReadFile err:", err)
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(clientKeyPath)
	if err != nil {
		log.Println("ReadFile err:", err)
		return nil, err
	}

	return NewETCDConn(caCrt, certPEMBlock, keyPEMBlock, endpoints)
}

func NewETCDConn(caCrt, clientCrt, clientKey []byte, endpoints []string) (*ETCDConn, error) {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.X509KeyPair(clientCrt, clientKey)
	if err != nil {
		log.Println("Loadx509keypair err:", err)
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
			//InsecureSkipVerify: true,
		},
	}

	return &ETCDConn{
		transport: transport,
		endpoints: endpoints,
	}, nil
}

func (p *ETCDConn) API(auth Auth) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: p.endpoints, // []string{"https://10.14.87.171:2379"},
		Transport: p.transport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
		Username:                auth.UserName, // "root",
		Password:                auth.Password, // "123456",
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client.NewKeysAPI(c), nil
}
