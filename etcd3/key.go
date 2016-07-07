package etcd3

import (
	"log"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"time"
	"strconv"
	"errors"
)

var requestTimeout = 20 * time.Second

func Demo(key string, auth3 *Auth3){
	cli, err := Get().API(auth3)
	if err != nil {
		log.Fatal(err)
	}

	kvc := clientv3.NewKV(cli)

	log.Println("init key ")
	if err := initKey(key, kvc); err != nil {
		log.Fatal(err)
	}

	for i:=0; i < 1000; i++ {
		value := strconv.Itoa(i)
		log.Println("begin ", value)

		if err := require(key,value, kvc); err != nil {
			log.Fatal(err)
		}

		if err := release(key,value, kvc); err != nil {
			log.Fatal(err)
		}
		
		log.Printf("end %s\n\n", value)
	}
}

func initKey(key string, kvc clientv3.KV) error {
	if _, err := kvc.Put(context.TODO(), key, ""); err != nil {
		return  err
	}
	return nil
}

func require(key, value string, kvc clientv3.KV) error {
	return cas(key, "", value, kvc)
}

func release(key, value string, kvc clientv3.KV) error {
	return cas(key, value, "", kvc)
}

func cas(key, value, newValue string, kvc clientv3.KV) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	tnxRes, err := kvc.Txn(ctx).
	If(clientv3.Compare(clientv3.Value(key), "=", value)).
	Then(clientv3.OpPut(key, newValue)).
	Else(clientv3.OpGet(key)).
	Commit()
	cancel()

	if err != nil {
		return err
	}

	if tnxRes.Succeeded {
		return nil
	}
	log.Println(string(tnxRes.Responses[0].GetResponseRange().Kvs[0].Value))
	return errors.New("release error")
}