package cherryDiscovery

import (
	"math/rand"
	"sync"

	cerr "github.com/cherry-game/cherry/error"
	cslice "github.com/cherry-game/cherry/extend/slice"
	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
	cprofile "github.com/cherry-game/cherry/profile"
)

// DiscoveryDefault 默认方式，通过读取profile文件的节点信息
//
// 该类型发现服务仅用于开发测试使用，直接读取profile.json->node配置
type DiscoveryDefault struct {
	sync.RWMutex
	memberMap        map[string]cfacade.IMember // key:nodeId,value:Member
	onAddListener    []cfacade.MemberListener
	onRemoveListener []cfacade.MemberListener
}

func (n *DiscoveryDefault) PreInit() {
	n.memberMap = map[string]cfacade.IMember{}
}

func (n *DiscoveryDefault) Load(_ cfacade.IApplication) {
	// load node info from profile file
	nodeConfig := cprofile.GetConfig("node")
	if nodeConfig == nil {
		clog.Error("`node` property not found in profile file.")
		return
	}

	for _, nodeType := range nodeConfig.Keys() {
		typeJson := nodeConfig.GetConfig(nodeType)
		for i := 0; i < typeJson.Size(); i++ {
			item := typeJson.GetConfig(i)

			nodeId := item.GetString("node_id")
			if nodeId == "" {
				clog.Errorf("nodeId is empty in nodeType = %s", nodeType)
				break
			}

			if _, found := n.GetMember(nodeId); found {
				clog.Errorf("nodeType = %s, nodeId = %s, duplicate nodeId", nodeType, nodeId)
				break
			}

			member := &cproto.Member{
				NodeId:   nodeId,
				NodeType: nodeType,
				Address:  item.GetString("rpc_address"),
				Settings: make(map[string]string),
			}

			settings := item.GetConfig("__settings__")
			for _, key := range settings.Keys() {
				member.Settings[key] = settings.GetString(key)
			}

			n.memberMap[member.NodeId] = member
		}
	}
}

func (n *DiscoveryDefault) Name() string {
	return "default"
}

func (n *DiscoveryDefault) Map() map[string]cfacade.IMember {
	return n.memberMap
}

func (n *DiscoveryDefault) ListByType(nodeType string, filterNodeId ...string) []cfacade.IMember {
	var list []cfacade.IMember

	for _, member := range n.memberMap {
		if member.GetNodeType() == nodeType {
			if _, ok := cslice.StringIn(member.GetNodeId(), filterNodeId); !ok {
				list = append(list, member)
			}
		}
	}

	return list
}

func (n *DiscoveryDefault) Random(nodeType string) (cfacade.IMember, bool) {
	memberList := n.ListByType(nodeType)
	memberLen := len(memberList)

	if memberLen < 1 {
		return nil, false
	}

	if memberLen == 1 {
		return memberList[0], true
	}

	return memberList[rand.Intn(len(memberList))], true
}

func (n *DiscoveryDefault) GetType(nodeId string) (nodeType string, err error) {
	member, found := n.GetMember(nodeId)
	if !found {
		return "", cerr.Errorf("nodeId = %s not found.", nodeId)
	}
	return member.GetNodeType(), nil
}

func (n *DiscoveryDefault) GetMember(nodeId string) (cfacade.IMember, bool) {
	if nodeId == "" {
		return nil, false
	}

	member, found := n.memberMap[nodeId]
	return member, found
}

func (n *DiscoveryDefault) AddMember(member cfacade.IMember) {
	defer n.Unlock()
	n.Lock()

	if _, found := n.GetMember(member.GetNodeId()); found {
		clog.Warnf("duplicate nodeId. [nodeType = %s], [nodeId = %s], [address = %s]",
			member.GetNodeType(),
			member.GetNodeId(),
			member.GetAddress(),
		)
		return
	}

	n.memberMap[member.GetNodeId()] = member

	for _, listener := range n.onAddListener {
		listener(member)
	}

	clog.Debugf("addMember new member. [member = %s]", member)
}

func (n *DiscoveryDefault) RemoveMember(nodeId string) {
	defer n.Unlock()
	n.Lock()

	member, found := n.GetMember(nodeId)
	if !found {
		return
	}

	delete(n.memberMap, member.GetNodeId())
	clog.Debugf("remove member. [member = %s]", member)

	for _, listener := range n.onRemoveListener {
		listener(member)
	}
}

func (n *DiscoveryDefault) OnAddMember(listener cfacade.MemberListener) {
	if listener == nil {
		return
	}
	n.onAddListener = append(n.onAddListener, listener)
}

func (n *DiscoveryDefault) OnRemoveMember(listener cfacade.MemberListener) {
	if listener == nil {
		return
	}
	n.onRemoveListener = append(n.onRemoveListener, listener)
}

func (n *DiscoveryDefault) Stop() {

}
