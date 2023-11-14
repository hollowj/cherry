package cherryProfile

import (
	"fmt"

	cerr "github.com/cherry-game/cherry/error"
	cfacade "github.com/cherry-game/cherry/facade"
)

// Node node info
type Node struct {
	nodeId     string
	nodeType   string
	address    string
	rpcAddress string
	settings   cfacade.ProfileCfg
	enabled    bool
}

func (n *Node) NodeId() string {
	return n.nodeId
}

func (n *Node) NodeType() string {
	return n.nodeType
}

func (n *Node) Address() string {
	return n.address
}

func (n *Node) RpcAddress() string {
	return n.rpcAddress
}

func (n *Node) Settings() cfacade.ProfileCfg {
	return n.settings
}

func (n *Node) Enabled() bool {
	return n.enabled
}

const stringFormat = "nodeId = %s, nodeType = %s, address = %s, rpcAddress = %s, enabled = %v"

func (n *Node) String() string {
	return fmt.Sprintf(stringFormat,
		n.nodeId,
		n.nodeType,
		n.address,
		n.rpcAddress,
		n.enabled,
	)
}

func GetNodeWithConfig(config *Config, nodeId string) (cfacade.INode, error) {
	nodeConfig := config.GetConfig("node")

	for _, nodeType := range nodeConfig.Keys() {
		typeJson := nodeConfig.GetConfig(nodeType)
		for i := 0; i < typeJson.Size(); i++ {
			item := typeJson.GetConfig(i)

			if nodeId != item.GetString("node_id") {
				continue
			}

			node := &Node{
				nodeId:     nodeId,
				nodeType:   nodeType,
				address:    item.GetString("address"),
				rpcAddress: item.GetString("rpc_address"),
				settings:   item.GetConfig("__settings__"),
				enabled:    item.GetBool("enabled"),
			}

			return node, nil
		}
	}

	return nil, cerr.Errorf("nodeId = %s not found.", nodeId)
}

func LoadNode(nodeId string) (cfacade.INode, error) {
	return GetNodeWithConfig(cfg.jsonConfig, nodeId)
}
