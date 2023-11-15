package cherryCluster

import (
	cfacade "github.com/cherry-game/cherry/facade"
	cherryRpcxCluster "github.com/cherry-game/cherry/net/cluster/rpcx_cluster"
)

const (
	Name = "cluster_component"
)

type Component struct {
	cfacade.Component
	cfacade.ICluster
}

func New() *Component {
	return &Component{}
}

func (c *Component) Name() string {
	return Name
}

func (c *Component) Init() {
	c.ICluster = c.loadCluster()
	c.ICluster.Init()
}

func (c *Component) OnStop() {
	c.ICluster.Stop()
}

func (c *Component) loadCluster() cfacade.ICluster {
	return cherryRpcxCluster.New(c.App())
}
