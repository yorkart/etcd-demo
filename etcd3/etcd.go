package etcd3

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	dialTimeout    = 5 * time.Second
)

type Auth3 struct {
	UserName string
	Password string
}

type ETCDConn struct {
//	transport *http.Transport
	t *tls.Config
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

	t := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cliCrt},
		//InsecureSkipVerify: true,
	}

	return &ETCDConn{
		t: t,
		endpoints: endpoints,
	}, nil
}

func (p *ETCDConn) API(auth *Auth3) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints: p.endpoints,
		TLS: p.t,

		// set timeout per request to fail fast when the target endpoint is unavailable
		DialTimeout: dialTimeout,
	}
	if auth != nil {
		cfg.Username = auth.UserName
		cfg.Password = auth.Password
	}

	return clientv3.New(cfg)
}

func (p *ETCDConn) AuthAPI(auth *Auth3) (clientv3.Auth,error) {
	cli, err := p.API(auth)
	if err != nil {
		return nil, err
	}

	return clientv3.NewAuth(cli),nil
}