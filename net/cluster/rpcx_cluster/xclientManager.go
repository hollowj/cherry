package cherryRpcxCluster

import (
	"sync"

	clog "github.com/cherry-game/cherry/logger"
	client2 "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
)

type XClientManager struct {
	xclientMap map[string]client.XClient
	mutex      sync.RWMutex
}

var xclientManagerOnce sync.Once
var xclientManagerIns *XClientManager

func GetXClientManager() *XClientManager {
	xclientManagerOnce.Do(func() {
		xclientManagerIns = &XClientManager{xclientMap: make(map[string]client.XClient)}
	})
	return xclientManagerIns
}

func (m *XClientManager) GetXClient(servicePath string) (client.XClient, error) {
	m.mutex.RLock()
	if xClient, ok := m.xclientMap[servicePath]; ok {
		m.mutex.RUnlock()
		return xClient, nil
	}
	m.mutex.RUnlock()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	clusterConfig := getClusterConfig()
	end_points := clusterConfig.GetString("end_points")
	prefix := clusterConfig.GetString("prefix")
	d, err := client2.NewConsulDiscovery(prefix, servicePath, []string{end_points}, nil)
	if err != nil {
		clog.Error(err)
		return nil, err
	}
	xclient := client.NewXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	m.xclientMap[servicePath] = xclient
	return xclient, nil
}
func (m *XClientManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, xClient := range m.xclientMap {
		err := xClient.Close()
		if err != nil {
			clog.Error(err)
		}
	}
	m.xclientMap = map[string]client.XClient{}
}
