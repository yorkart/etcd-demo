package main
import (
	"strconv"
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"etcd-demo/etcd"
)

var auth = etcd.Auth{UserName: "root", Password: "123456"}
var key = "/demo/a"

func main() {
	etcd.CreateConn()
	keyAPI, err := etcd.Get().API(auth);
	if err != nil {
		log.Println("get key api error: ", err)
		return
	}

	for i:=0;i < 1000; i++ {
		value := strconv.Itoa(i)

		log.Println("=> atomic set value: ", value)
		_, err = keyAPI.Set(context.Background(),key, value, &client.SetOptions{PrevExist: client.PrevNoExist})
		if err != nil {
			log.Println("set key error: ", err)
			return
		}

		time.Sleep(time.Second)

		log.Println("<= delete value: ", value)
		_, err = keyAPI.Delete(context.Background(), key, nil)
		if err != nil {
			log.Println("delete key error: ", err)
			return
		}

		time.Sleep(time.Second)

		log.Println("finish ", value)
	}

}
