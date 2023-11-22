package cherryDiscovery

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	clog "github.com/cherry-game/cherry/logger"
	"github.com/hashicorp/consul/api"
)

func TestConsul(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	key := "k"
	value := "v"
	sessionId, _, err2 := client.Session().Create(&api.SessionEntry{Name: key, TTL: "10s", Behavior: api.SessionBehaviorDelete}, nil)
	if err2 != nil {
		panic(err)

	}
	lock, err := client.LockOpts(&api.LockOptions{Key: key, Session: sessionId})
	if err != nil {
		panic(err)
	}
	clog.Info(lock)
	_, err = lock.Lock(make(chan struct{}, 1))
	//defer lock.Unlock()
	if err != nil {
		panic(err)

	}
	pair := &api.KVPair{Key: key, Value: []byte(value), Session: sessionId, Flags: api.LockFlagValue}
	_, err = client.KV().Put(pair, nil)
	if err != nil {
		panic(err)
	}
	clog.Info(sessionId)
}

func TestConsul2(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	key := "k"
	value := "vv"
	sessionId, _, err2 := client.Session().Create(&api.SessionEntry{Name: key, TTL: "10s", Behavior: api.SessionBehaviorDelete}, nil)
	if err2 != nil {
		panic(err)

	}
	lock, err := client.LockOpts(&api.LockOptions{Key: key, Session: sessionId})
	if err != nil {
		panic(err)
	}
	clog.Info(lock)
	_, err = lock.Lock(make(chan struct{}, 1))
	//defer lock.Unlock()
	if err != nil {
		panic(err)

	}
	pair := &api.KVPair{Key: key, Value: []byte(value), Session: sessionId, Flags: api.LockFlagValue}
	_, err = client.KV().Put(pair, nil)
	if err != nil {
		panic(err)
	}
	clog.Info(sessionId)
}

func TestConsul3(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	key := "k"
	value := "vvv"

	pair := &api.KVPair{Key: key, Value: []byte(value)}
	_, err = client.KV().Put(pair, nil)
	if err != nil {
		panic(err)
	}
}
func Test4(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println(len(ch))

			}
		}
	}()
	time.Sleep(time.Second * 5)
	close(ch)
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case s := <-sg:
		clog.Infof("receive shutdown signal = %v.", s)
	}
}
func TestHttp(t *testing.T) {
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		clog.Info("ping")
		writer.Write([]byte("pong"))
	})
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hi"))
	})
	http.ListenAndServe(":8080", nil)
}
func TestService(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	reg := &api.AgentServiceRegistration{
		ID:      "hello_service1",
		Name:    "hello",
		Port:    8080,
		Address: "127.0.0.1",
		Tags:    []string{"hello"},
	}
	reg.Check = &api.AgentServiceCheck{HTTP: fmt.Sprintf("http://%s:%d/ping1", reg.Address, reg.Port), Interval: "5s", Timeout: "3s", DeregisterCriticalServiceAfter: "10s"}
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		panic(err)
	}
	//sg := make(chan os.Signal, 1)
	//signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	//
	//select {
	//case s := <-sg:
	//	clog.Infof("receive shutdown signal = %v.", s)
	//}
}
func TestServiceFind(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	m, err := client.Agent().ServicesWithFilter("hello in Tags")
	if err != nil {
		panic(err)
	}
	for s, service := range m {
		clog.Infof("%v=%+v", s, service)
	}

	//sg := make(chan os.Signal, 1)
	//signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	//
	//select {
	//case s := <-sg:
	//	clog.Infof("receive shutdown signal = %v.", s)
	//}
}
func TestHealthyServiceFind(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}

	entries, _, err := client.Health().Service("hello", "hello", true, nil)
	if err != nil {
		return
	}
	for s, service := range entries {
		clog.Infof("%v=%+v", s, service.Service)
	}
	//sg := make(chan os.Signal, 1)
	//signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	//
	//select {
	//case s := <-sg:
	//	clog.Infof("receive shutdown signal = %v.", s)
	//}
}
func TestSession(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		clog.Fatal(err)
	}
	sessionId, _, err := client.Session().Create(&api.SessionEntry{TTL: "10s", Behavior: api.SessionBehaviorDelete}, nil)
	if err != nil {
		panic(err)
	}
	k := "h"
	v := "wj"
	_, _, err = client.KV().Acquire(&api.KVPair{Key: k, Value: []byte(v), Session: sessionId}, nil)
	if err != nil {
		panic(err)
	}
	//client.LockOpts(&api.LockOptions{Key: k,})
}
