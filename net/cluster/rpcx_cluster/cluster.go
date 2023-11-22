package cherryRpcxCluster

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	cherryNet "github.com/cherry-game/cherry/extend/net"
	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap/zapcore"
)

type (
	Cluster struct {
		app       cfacade.IApplication
		rpcServer *server.Server
		service   *ClusterService
		rpcClient *ClusterClient
		rpcAddr   string
	}

	OptionFunc func(o *Cluster)
)

func LocalIPWithAutoPortForGame(min int32, max int32) string {
	ip := cherryNet.LocalIPV4()
	rand.Seed(time.Now().UnixNano())
	for {
		site := fmt.Sprintf("%s:%d", ip, min+rand.Int31n(max-min))
		listen, err := net.Listen("tcp", site)
		if err == nil {
			err = listen.Close()
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			return site
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func New(app cfacade.IApplication, options ...OptionFunc) cfacade.ICluster {
	addr := LocalIPWithAutoPortForGame(20000, 40000)
	s := server.NewServer()
	clusterClient := NewClusterClient(app)
	clusterService := NewClusterService(app)
	clusterConfig := getClusterConfig()
	consulAddr := clusterConfig.GetString("addr")
	prefix := clusterConfig.GetString("prefix")
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		ConsulServers:  []string{consulAddr},
		BasePath:       prefix,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Second * 10,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	s.Plugins.Add(r)
	err = s.RegisterName(getServiceName(app.NodeType(), app.NodeId()), clusterService, "")
	if err != nil {
		panic(err)
	}
	cluster := &Cluster{
		app:       app,
		rpcServer: s,
		service:   clusterService,
		rpcClient: clusterClient,
		rpcAddr:   addr,
	}

	for _, option := range options {
		option(cluster)
	}
	return cluster
}

func (p *Cluster) Init() {
	go func() {
		err := p.rpcServer.Serve("tcp", p.rpcAddr)
		if err != nil {
			panic(err)
		}
	}()
	clog.Info("rpcx cluster execute OnInit().")
}

func (p *Cluster) Stop() {
	err := p.rpcServer.UnregisterAll()
	if err != nil {
		clog.Error(err)
	}
	err = p.rpcServer.Shutdown(context.TODO())
	if err != nil {
		clog.Error(err)
	}
	xclientManagerIns.Close()
	clog.Info("rpcx cluster execute OnStop().")
}

func (p *Cluster) PublishLocal(nodeId string, request *cproto.ClusterPacket) error {
	if clog.PrintLevel(zapcore.DebugLevel) {
		clog.Debugf("[PublishLocal] [nodeId = %s, %s]",
			nodeId,
			request.PrintLog(),
		)
	}
	err := p.rpcClient.PublishLocal(nodeId, request)
	if err != nil {
		return err
	}

	return err
}

func (p *Cluster) PublishRemote(nodeId string, request *cproto.ClusterPacket) error {
	if clog.PrintLevel(zapcore.DebugLevel) {
		clog.Debugf("[PublishRemote] [nodeId = %s, %s]",
			nodeId,
			request.PrintLog(),
		)
	}
	err := p.rpcClient.PublishRemote(nodeId, request)

	return err
}

func (p *Cluster) RequestRemote(nodeId string, request *cproto.ClusterPacket, timeout ...time.Duration) cproto.Response {
	if clog.PrintLevel(zapcore.DebugLevel) {
		clog.Debugf("[RequestRemote] [nodeId = %s, %s]",
			nodeId,
			request.PrintLog(),
		)
	}
	rsp := p.rpcClient.RequestRemote(nodeId, request)
	return rsp
}
