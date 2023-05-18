package cherry

import (
	"context"
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cagent "github.com/cherry-game/cherry/net/agent"
	ccluster "github.com/cherry-game/cherry/net/cluster"
	ccommand "github.com/cherry-game/cherry/net/command"
	ccontext "github.com/cherry-game/cherry/net/context"
	chandler "github.com/cherry-game/cherry/net/handler"
	cmsg "github.com/cherry-game/cherry/net/message"
	cpacket "github.com/cherry-game/cherry/net/packet"
	cproto "github.com/cherry-game/cherry/net/proto"
	crouter "github.com/cherry-game/cherry/net/router"
	csession "github.com/cherry-game/cherry/net/session"
)

var (
	_thisApp    *Application
	_components []cfacade.IComponent
)

var (
	_commands         = make(map[cpacket.Type]ccommand.ICommand)
	_handshakeData    = make(map[string]interface{})
	_heartbeat        = 60 * time.Second
	_connectors       []cfacade.IConnector
	_clusterComponent *ccluster.Component
)

var (
	_handlerOpts      []chandler.Option
	_handlerGroups    []*chandler.HandlerGroup
	_handlerComponent *chandler.Component
)

func App() *Application {
	return _thisApp
}

func Configure(profilePath, profileName, nodeId string) cfacade.IApplication {
	_thisApp = NewApp(profilePath, profileName, nodeId)
	return _thisApp
}

func Run(isFrontend bool, nodeMode NodeMode) {
	if _thisApp == nil {
		panic("please call the configure function first.")
	}

	if _thisApp.Running() {
		return
	}

	_thisApp.isFrontend = isFrontend
	_thisApp.nodeMode = nodeMode

	initHandler()
	initCluster()
	initConnector()
	initComponent()

	_thisApp.Startup()
}

func initHandler() {
	// register handler component
	_handlerComponent = chandler.NewComponent(_handlerOpts...)

	for _, group := range _handlerGroups {
		_handlerComponent.Register(group)
	}

	// add handler component
	_thisApp.Register(_handlerComponent)
}

func initComponent() {
	_thisApp.Register(_components...)
}

func initCluster() {
	if _thisApp.NodeMode() == Cluster {
		// register cluster component
		_clusterComponent = ccluster.NewComponent()
		_thisApp.Register(_clusterComponent)
	}
}

func initConnector() {
	if _thisApp.isFrontend == false {
		return
	}

	if len(_connectors) < 1 {
		panic("please call the cherry.RegisterConnector() method add IConnector.")
	}

	initCommand()
	initOnSession()

	for _, connector := range _connectors {
		// default setting
		if connector.IsSetListener() == false {
			connector.OnConnectListener(func(conn cfacade.INetConn) {
				// create agent
				agent := cagent.NewAgent(_thisApp, conn, &cagent.Options{
					Heartbeat: _heartbeat,
					Commands:  _commands,
				})

				// create new session
				csession.Create(csession.NextSID(), _thisApp.NodeId(), agent, func(session *csession.Session) {
					agent.SetSession(session)
				})
				// run agent
				agent.Run()
			})
		}

		RegisterComponent(connector)
	}
}

func initOnSession() {
	csession.AddOnCreateListener(func(session *csession.Session) (next bool) {
		session.Debugf("session create. [sid = %s, address = %s]",
			session.SID(),
			session.RemoteAddress(),
		)
		return true
	})

	csession.AddOnCloseListener(func(session *csession.Session) (next bool) {
		session.Debugf("session closed. [sid = %s, address = %s]",
			session.SID(),
			session.RemoteAddress(),
		)

		return true
	})
}

func initCommand() {
	if _, found := _commands[cpacket.Handshake]; !found {
		if len(_handshakeData) < 1 {
			_handshakeData["heartbeat"] = _heartbeat.Seconds()
			_handshakeData["dict"] = cmsg.GetDictionary()
			_handshakeData["serializer"] = _thisApp.ISerializer.Name()
		}

		handshakeCommand := ccommand.NewHandshake(_thisApp, _handshakeData)
		RegisterCommand(handshakeCommand)
	}

	if _, found := _commands[cpacket.HandshakeAck]; !found {
		handshakeAckCommand := ccommand.NewHandshakeACK()
		RegisterCommand(handshakeAckCommand)
	}

	if _, found := _commands[cpacket.Heartbeat]; !found {
		heartbeatCommand := ccommand.NewHeartbeat(_thisApp)
		RegisterCommand(heartbeatCommand)
	}

	if _, found := _commands[cpacket.Data]; !found {
		// connector forward message
		handDataCommand := ccommand.NewData(
			_thisApp,
			_handlerComponent.ProcessLocal,
			forwardLocal,
		)
		RegisterCommand(handDataCommand)
	}
}

// ForwardLocal forward message to backend node
func forwardLocal(ctx context.Context, session *csession.Session, msg *cmsg.Message) {
	member, err := crouter.Route(context.Background(), msg.RouteInfo(), session)
	if member == nil || err != nil {
		session.Warnf("forward node fail. [session = %v, msg = %s, error = %s]",
			session,
			msg,
			err,
		)
		return
	}

	request := buildRequest(ctx, session, msg)
	defer request.Recycle()

	if err = _clusterComponent.PublishLocal(member.GetNodeId(), request); err != nil {
		session.Warnf("publish local fail. [nodeId = %s, route = %s, error = %s]",
			member.GetNodeId(),
			msg.Route,
			err,
		)
	}
}

func SetSerializer(serializer cfacade.ISerializer) {
	_thisApp.SetSerializer(serializer)
}

func SetPacketCodec(codec cfacade.IPacketCodec) {
	_thisApp.SetPacketCodec(codec)
}

func SetHeartbeat(t time.Duration) {
	if t.Seconds() < 1 {
		t = 60 * time.Second
	}
	_heartbeat = t
}

func SetHandshake(key string, value interface{}) {
	_handshakeData[key] = value
}

func SetDictionary(dict map[string]uint16) {
	cmsg.SetDictionary(dict)
}

func SetDataCompression(compression bool) {
	cmsg.SetDataCompression(compression)
}

func SetOnShutdown(fn ...func()) {
	_thisApp.OnShutdown(fn...)
}

func SetHandlerOptions(opts ...chandler.Option) {
	_handlerOpts = append(_handlerOpts, opts...)
}

func RegisterHandler(handler ...cfacade.IHandler) {
	handlerGroup := chandler.NewGroupWithHandler(handler...)
	_handlerGroups = append(_handlerGroups, handlerGroup)
}

func RegisterHandlerGroup(group ...*chandler.HandlerGroup) {
	_handlerGroups = append(_handlerGroups, group...)
}

func RegisterComponent(component ...cfacade.IComponent) {
	_components = append(_components, component...)
}

func RegisterConnector(connector cfacade.IConnector) {
	_connectors = append(_connectors, connector)
}

func RegisterCommand(command ccommand.ICommand) {
	_commands[command.PacketType()] = command
}

func AddNodeRouter(nodeType string, routingFunc crouter.RoutingFunc) {
	crouter.AddRoute(nodeType, routingFunc)
}

func GetConnectors() []cfacade.IConnector {
	return _connectors
}

func PostEvent(event cfacade.IEvent) {
	if _handlerComponent == nil {
		clog.Warnf("post event fail. handler component is nil.")
		return
	}

	_handlerComponent.PostEvent(event)
}

func buildRequest(ctx context.Context, session *csession.Session, msg *cmsg.Message) *cproto.Request {
	request := cproto.GetRequest()
	request.Sid = session.SID()
	request.Uid = session.UID()
	request.FrontendId = session.FrontendId()
	request.Setting = session.Data()
	request.Context = ccontext.Encode(ctx)

	request.MsgType = int32(msg.Type)
	request.MsgId = uint32(msg.ID)
	request.Route = msg.Route
	request.Data = msg.Data
	request.IsError = false

	return request
}
