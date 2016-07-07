package main
import (
	"log"

	"etcd-demo/etcd"
	"etcd-demo/etcd3"
	"time"
)

var auth = etcd.Auth{UserName: "web", Password: "warden@web"}
var auth3 = &etcd3.Auth3{UserName: "web", Password: "warden@web"}

var key = "/warden/key"

func main() {
	v3()

	time.Sleep(3*time.Second)
}

func v2() {
	if err := etcd.CreateConn() ; err != nil {
		log.Fatal(err)
	}
	etcd.Demo(key, auth)
}

func v3() {
	if err := etcd3.CreateConn() ; err != nil {
		log.Fatal(err)
	}
	etcd3.Demo(key, auth3)
}