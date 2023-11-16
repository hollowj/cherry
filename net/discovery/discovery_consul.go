package cherryDiscovery

import (
	"fmt"
	"strings"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
	cprofile "github.com/cherry-game/cherry/profile"
	jsoniter "github.com/json-iterator/go"
	"github.com/rpcxio/libkv/store"
	"github.com/rpcxio/libkv/store/consul"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

var (
	ConsulKeyPrefix         = "node/"
	ConsulRegisterKeyFormat = ConsulKeyPrefix + "%s"
)

// Consul etcd方式发现服务
type Consul struct {
	app cfacade.IApplication
	DiscoveryDefault
	prefix string
	ttl    int64
	kv     store.Store
	addr   string
}

func NewConsul() *Consul {
	return &Consul{}
}

func (p *Consul) Name() string {
	return "consul"
}

func (p *Consul) Load(app cfacade.IApplication) {
	p.DiscoveryDefault.PreInit()
	p.app = app
	p.ttl = 10

	clusterConfig := cprofile.GetConfig("cluster").GetConfig(p.Name())
	if clusterConfig == nil {
		clog.Fatalf("Consul config not found.")
		return
	}

	p.loadConfig(clusterConfig)
	p.init()
	p.register()
	p.watch()

	clog.Infof("[Consul] init complete! [addr = %v] ", p.addr)
}

func (p *Consul) OnStop() {
	key := fmt.Sprintf(ConsulRegisterKeyFormat, p.app.NodeId())
	err := p.kv.Delete(key)
	clog.Infof("Consul stopping! err = %v", err)

}

func (p *Consul) loadConfig(config cfacade.ProfileCfg) {

	p.addr = config.GetString("addr")
	p.ttl = config.GetInt64("ttl", 5)
	p.prefix = config.GetString("prefix", "cherry")
}

func (p *Consul) init() {
	kv, err := consul.New([]string{p.addr}, nil)
	if err != nil {
		panic(err)
	}
	p.kv = kv

}

func (p *Consul) register() {
	registerMember := &cproto.Member{
		NodeId:   p.app.NodeId(),
		NodeType: p.app.NodeType(),
		Address:  p.app.RpcAddress(),
		Settings: make(map[string]string),
	}

	jsonString, err := jsoniter.MarshalToString(registerMember)
	if err != nil {
		clog.Fatal(err)
		return
	}

	key := fmt.Sprintf(ConsulRegisterKeyFormat, p.app.NodeId())

	err = p.kv.Put(key, []byte(jsonString), nil)
	if err != nil {
		clog.Fatal(err)
		return
	}
}

func (p *Consul) watch() {
	resp, err := p.kv.List(ConsulKeyPrefix)
	if err != nil {
		clog.Fatal(err)
		return
	}

	for _, ev := range resp {
		p.addMember(ev.Value)
	}
	stopCh := make(<-chan struct{})

	watchChan, err := p.kv.WatchTree(ConsulKeyPrefix, stopCh)
	go func() {
		for rsp := range watchChan {
			p.updateMembers(rsp)
		}
	}()
}

func (p *Consul) addMember(data []byte) {
	member := &cproto.Member{}
	err := jsoniter.Unmarshal(data, member)
	if err != nil {
		return
	}

	p.AddMember(member)
}

func (p *Consul) removeMember(kv *mvccpb.KeyValue) {
	key := string(kv.Key)
	nodeId := strings.ReplaceAll(key, ConsulKeyPrefix, "")
	if nodeId == "" {
		clog.Warn("remove member nodeId is empty!")
	}

	p.RemoveMember(nodeId)
}
func (n *DiscoveryDefault) updateMembers(kvPairs []*store.KVPair) {
	memberMap := make(map[string]cfacade.IMember) // key:nodeId,value:Member

	for _, pair := range kvPairs {
		member := &cproto.Member{}
		err := jsoniter.Unmarshal(pair.Value, member)
		if err != nil {
			return
		}
		memberMap[member.GetNodeId()] = member
	}
	defer n.Unlock()
	n.Lock()
	n.memberMap = memberMap

	clog.Debugf("addMember new member. [member = %s]", memberMap)
}
