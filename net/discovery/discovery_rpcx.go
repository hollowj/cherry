package cherryDiscovery

import (
	"strings"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/rpcxio/libkv/store"
	"github.com/rpcxio/libkv/store/consul"
)

var (
	ConsulRpcxKeyPrefix = "cherry/nodes/"
)

// ConsulRpcx consul方式发现服务
type ConsulRpcx struct {
	app cfacade.IApplication
	DiscoveryDefault
	prefix string
	ttl    int64
	kv     store.Store
	addr   string
}

func NewRpcxConsul() *ConsulRpcx {
	return &ConsulRpcx{}
}

func (p *ConsulRpcx) Name() string {
	return "consul"
}

func (p *ConsulRpcx) Load(app cfacade.IApplication) {
	p.DiscoveryDefault.PreInit()
	p.app = app
	p.ttl = 10

	clusterConfig := cprofile.GetConfig("cluster").GetConfig(p.Name())
	if clusterConfig == nil {
		clog.Fatalf("ConsulRpcx config not found.")
		return
	}

	p.loadConfig(clusterConfig)
	p.init()
	p.watch()

	clog.Infof("[ConsulRpcx] init complete! [addr = %v] ", p.addr)
}

func (p *ConsulRpcx) OnStop() {

	clog.Infof("ConsulRpcx stopping! ")

}

func (p *ConsulRpcx) loadConfig(config cfacade.ProfileCfg) {

	p.addr = config.GetString("addr")
	p.ttl = config.GetInt64("ttl", 5)
	p.prefix = config.GetString("prefix", "cherry")
}

func (p *ConsulRpcx) init() {
	kv, err := consul.New([]string{p.addr}, nil)
	if err != nil {
		panic(err)
	}
	p.kv = kv

}

func (p *ConsulRpcx) watch() {
	resp, err := p.kv.List(ConsulRpcxKeyPrefix)
	if err != nil {
		clog.Fatal(err)
		return
	}
	p.updateMembers(resp)
	stopCh := make(chan struct{})

	watchChan, err := p.kv.WatchTree(ConsulRpcxKeyPrefix, stopCh)
	go func() {
		for {
			select {
			case rsp := <-watchChan:
				if rsp == nil {
					clog.Error("ConsulRpcx maybe is down")
					return
				}
				p.updateMembers(rsp)
			case die := <-p.app.DieChan():
				if die {
					close(stopCh)
				}
			}

		}
	}()

}

func (n *ConsulRpcx) updateMembers(kvPairs []*store.KVPair) {
	memberMap := make(map[string]cfacade.IMember) // key:nodeId,value:Member

	for _, pair := range kvPairs {
		if strings.Count(pair.Key, "/") == 2 {
			typeStr := pair.Key[strings.LastIndex(pair.Key, "/")+1:]
			strArr := strings.Split(typeStr, "@")
			if len(strArr) == 2 {
				member := &cproto.Member{
					NodeId:   strArr[1],
					NodeType: strArr[0],
				}
				memberMap[member.GetNodeId()] = member
			}
		}

	}
	n.Lock()
	defer n.Unlock()
	if len(memberMap) != len(n.memberMap) {
		clog.Debugf("updateMembers new member. [member = %s]", memberMap)
	}
	n.memberMap = memberMap

}
