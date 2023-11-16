package cherryRpcxCluster

import (
	"context"
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cproto "github.com/cherry-game/cherry/net/proto"
)

type ClusterService struct {
	app cfacade.IApplication
}

func NewClusterService(app cfacade.IApplication) *ClusterService {
	return &ClusterService{app: app}
}
func (p *ClusterService) PublishLocal(ctx context.Context, request *cproto.ClusterPacket, response *cproto.ClusterPacket) error {
	message := cfacade.GetMessage()
	message.BuildTime = request.BuildTime
	message.Source = request.SourcePath
	message.Target = request.TargetPath
	message.FuncName = request.FuncName
	message.IsCluster = true
	message.Session = request.Session
	message.Args = request.ArgBytes

	p.app.ActorSystem().PostLocal(message)

	return nil
}

func (p *ClusterService) PublishRemote(ctx context.Context, request *cproto.ClusterPacket, response *cproto.ClusterPacket) error {

	message := cfacade.GetMessage()
	message.BuildTime = request.BuildTime
	message.Source = request.SourcePath
	message.Target = request.TargetPath
	message.FuncName = request.FuncName
	if request.ArgBytes != nil {
		message.Args = request.ArgBytes
	}

	message.IsCluster = true
	p.app.ActorSystem().PostRemote(message)
	return nil
}

func (p *ClusterService) RequestRemote(ctx context.Context, request *cproto.ClusterPacket, response *cproto.ClusterPacket) error {
	message := cfacade.GetMessage()
	message.BuildTime = request.BuildTime
	message.Source = request.SourcePath
	message.Target = request.TargetPath
	message.FuncName = request.FuncName
	if request.ArgBytes != nil {
		message.Args = request.ArgBytes
	}

	message.IsCluster = true
	res := newResponse()
	message.ClusterReply = res
	p.app.ActorSystem().PostRemote(message)
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*3)
	defer cancelFunc()
	select {
	case data := <-res.dataChan:
		response.ArgBytes = data
	case <-timeout.Done():
		clog.Warnf("RequestRemote timeout:%s", request.PrintLog())
	}
	return nil
}

type response struct {
	dataChan chan []byte
}

func newResponse() *response {
	return &response{dataChan: make(chan []byte, 1)}
}
func (p *response) Respond(data []byte) error {
	p.dataChan <- data
	return nil
}
