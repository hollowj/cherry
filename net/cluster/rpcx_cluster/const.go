package cherryRpcxCluster

import (
	"fmt"

	cfacade "github.com/cherry-game/cherry/facade"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/spf13/viper"
)

const (
	serviceNameFormat = "nodes/%s@%s"
)

func getServiceName(nodeType string, nodeId string) string {
	return fmt.Sprintf(serviceNameFormat, nodeType, nodeId)
}

func getClusterConfig() cfacade.ProfileCfg {
	sMode := viper.GetString("cluster.discovery.mode")
	clusterConfig := cprofile.GetConfig("cluster").GetConfig(sMode)
	return clusterConfig
}
