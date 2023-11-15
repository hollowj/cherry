package cherryRpcxCluster

import (
	cfacade "github.com/cherry-game/cherry/facade"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/spf13/viper"
)

func getEtcdConfig() cfacade.ProfileCfg {
	sMode := viper.GetString("cluster.discovery.mode")
	clusterConfig := cprofile.GetConfig("cluster").GetConfig(sMode)
	return clusterConfig
}
