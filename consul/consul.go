package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"caton/xw.wang/utils"

	clog "github.com/cihub/seelog"
	consul_api "github.com/hashicorp/consul/api"
	consul_watch "github.com/hashicorp/consul/watch"
)

var (
	ConsulClient    *consul_api.Client
	KV              *consul_api.KV
	ConsulClientBak *consul_api.Client
	KVBak           *consul_api.KV

	CurrentClient int // 当前使用的consul客户端 1：ConsulClient 2:ConsulClientBak -1都不可用
)
var (
	chClientA chan string
	chClientB chan string
)

const (
	TIME_OUT  time.Duration = 2 * time.Second
	CLIENT_A                = 1
	CLIENT_B                = 2
	CLIENT_NO               = -1
)

func Put(key string, value string) (bool, error) {

	var err error
	chClientA = make(chan string, 1)
	chClientB = make(chan string, 1)

	go putA(key, value)
	//	go putB(key, value)

	select {
	case <-time.After(TIME_OUT):
		err = fmt.Errorf("consul time out.")
	case <-chClientA:
	case <-chClientB:
	}

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func putA(key string, value string) {
	var err error
	p := &consul_api.KVPair{Key: key, Value: []byte(value)}
	_, err = KV.Put(p, nil)
	if err != nil {
		//		clog.Error(err)
		return
	}
	chClientA <- ""
	return
}
func putB(key string, value string) {
	var err error
	p := &consul_api.KVPair{Key: key, Value: []byte(value)}
	_, err = KVBak.Put(p, nil)
	if err != nil {
		//		clog.Error(err)
		return
	}
	chClientB <- ""
	return
}

func Get(key string) (ret string, err error) {

	chClientA = make(chan string, 1)
	chClientB = make(chan string, 1)

	go getA(key)
	//	go getB(key)

	select {
	case <-time.After(TIME_OUT):
		err = fmt.Errorf("consul time out.")
		ret = ""
	case ret = <-chClientA:
	case ret = <-chClientB:
	}

	if err != nil {
		return "", err
	} else {
		return ret, nil
	}
}

func getA(key string) {
	var err error
	var pair *consul_api.KVPair
	pair, _, err = KV.Get(key, nil)
	if err != nil {
		clog.Error(err)
		return
	}
	chClientA <- string(pair.Value)
	return
}
func getB(key string) {
	var err error
	var pair *consul_api.KVPair
	pair, _, err = KVBak.Get(key, nil)
	if err != nil {
		clog.Error(err)
		return
	}
	chClientB <- string(pair.Value)
	return
}

func MustParse(q string) *consul_watch.WatchPlan {
	params := MakeParams(q)
	plan, err := consul_watch.Parse(params)
	if err != nil {
		clog.Error(err)
	}
	return plan
}
func MakeParams(s string) map[string]interface{} {
	var out map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	err := dec.Decode(&out)
	if err != nil {
		clog.Error(err)
	}
	return out
}

func InitClient(consulClientAddr string) {

	//	utils.SetEnv("CONSUL_HTTP_ADDR", common.ConsulClientAddrBak)
	//	ConsulClientBak, _ = consul_api.NewClient(consul_api.DefaultConfig())
	//	KVBak = ConsulClientBak.KV()

	utils.SetEnv("CONSUL_HTTP_ADDR", consulClientAddr)
	ConsulClient, _ = consul_api.NewClient(consul_api.DefaultConfig())
	KV = ConsulClient.KV()

}
