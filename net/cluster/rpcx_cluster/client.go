package cherryRpcxCluster

import (
	"context"
	"time"

	ccode "github.com/cherry-game/cherry/code"
	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
	"google.golang.org/protobuf/proto"
)

type ClusterClient struct {
	app cfacade.IApplication
}

func NewClusterClient(app cfacade.IApplication) *ClusterClient {
	return &ClusterClient{app: app}

}
func (p *ClusterClient) PublishLocal(nodeId string, request *cproto.ClusterPacket) error {
	defer request.Recycle()

	nodeType, err := p.app.Discovery().GetType(nodeId)
	if err != nil {
		clog.Debugf("[PublishLocal] get node type fail. [nodeId = %s, %s]",
			nodeId,
			request.PrintLog(),
		)
		return err
	}

	serviceName := getServiceName(nodeType, nodeId)
	xClient, err := GetXClientManager().GetXClient(serviceName)
	if err != nil {
		return err
	}
	reply := cproto.GetClusterPacket()
	err = xClient.Call(context.Background(), "PublishLocal", request, reply)
	if err != nil {
		clog.Errorf("failed to call: %v", err)
		return err
	}
	return err
}

func (p *ClusterClient) PublishRemote(nodeId string, request *cproto.ClusterPacket) error {
	defer request.Recycle()

	nodeType, err := p.app.Discovery().GetType(nodeId)
	if err != nil {
		clog.Debugf("[PublishRemote] Get node type fail. [nodeId = %s, %s, err = %v]",
			nodeId,
			request.PrintLog(),
			err,
		)
		return err
	}

	serviceName := getServiceName(nodeType, nodeId)
	xClient, err := GetXClientManager().GetXClient(serviceName)
	if err != nil {
		return err
	}
	reply := cproto.GetClusterPacket()
	err = xClient.Call(context.Background(), "PublishRemote", request, reply)
	if err != nil {
		clog.Errorf("failed to call: %v", err)
		return err
	}

	return err
}

func (p *ClusterClient) RequestRemote(nodeId string, request *cproto.ClusterPacket, timeout ...time.Duration) cproto.Response {
	defer request.Recycle()

	rsp := cproto.Response{}
	nodeType, err := p.app.Discovery().GetType(nodeId)
	if err != nil {
		clog.Debugf("[PublishRemote] Get node type fail. [nodeId = %s, %s, err = %v]",
			nodeId,
			request.PrintLog(),
			err,
		)

		rsp.Code = ccode.DiscoveryNotFoundNode
		return rsp
	}

	serviceName := getServiceName(nodeType, nodeId)
	xClient, err := GetXClientManager().GetXClient(serviceName)
	if err != nil {
		clog.Debugf("[GetXClient] Get node type fail. [nodeId = %s, %s, err = %v]",
			serviceName,
			request.PrintLog(),
			err,
		)
		rsp.Code = ccode.DiscoveryNotFoundNode
		return rsp
	}
	reply := cproto.GetClusterPacket()
	ctx := context.Background()
	if len(timeout) > 0 {
		withoutCtx, cancelFunc := context.WithTimeout(ctx, timeout[0])
		defer cancelFunc()
		ctx = withoutCtx
	}

	err = xClient.Call(ctx, "RequestRemote", request, reply)
	if err != nil {
		clog.Errorf("failed to call: %v %v", request.PrintLog(), err)
		rsp.Code = ccode.NodeRequestError
		return rsp
	}
	if err = proto.Unmarshal(reply.ArgBytes, &rsp); err != nil {
		clog.Warnf("[RequestRemote] unmarshal fail. [nodeId = %s, %s, rsp = %v, err = %v]",
			nodeId,
			request.PrintLog(),
			rsp,
			err,
		)

		rsp.Code = ccode.RPCUnmarshalError
		return rsp
	}

	return rsp
}
