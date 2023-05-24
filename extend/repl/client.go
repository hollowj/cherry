// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package repl

import (
	"encoding/json"

	cherryLogger "github.com/cherry-game/cherry/logger"
	cherryClient "github.com/cherry-game/cherry/net/client"
	cmessage "github.com/cherry-game/cherry/net/message"
	cherrySerializer "github.com/cherry-game/cherry/net/serializer"
)

// Client struct
type Client struct {
	*cherryClient.Client

	IncomingMsgChan chan *cmessage.Message
}

// ConnectedStatus returns the the current connection status
func (pc *Client) ConnectedStatus() bool {
	return pc.IsConnected()
}

// MsgChannel return the incoming message channel
func (pc *Client) MsgChannel() chan *cmessage.Message {
	return pc.IncomingMsgChan
}

// Return the basic structure for the Client struct.
func newClient() *Client {
	return &Client{
		Client:          cherryClient.New(),
		IncomingMsgChan: make(chan *cmessage.Message, 10),
	}
}

// NewWebsocket returns a new websocket client
func NewWebsocket(path string) *Client {

	wc := &Client{
		Client:          cherryClient.New(cherryClient.WithSerializer(cherrySerializer.NewJSON())),
		IncomingMsgChan: make(chan *cmessage.Message, 10),
	}
	wc.OnUnexpectedEvent(UnexpectedEventCb(wc))
	return wc
}

// UnexpectedEventCb returns a function to deal with un listened event
func UnexpectedEventCb(pc *Client) func(data interface{}) {
	return func(data interface{}) {
		push := data.(*cmessage.Message)
		pushStruct, err := routeMessage(push.Route)
		if err != nil {
			cherryLogger.Error(err.Error())
			return
		}

		err = pc.Serializer().Unmarshal(push.Data, pushStruct)
		if err != nil {
			cherryLogger.Error("unmarshal error data:%v ", push.Data)
			return
		}

		jsonData, err := json.Marshal(pushStruct)
		if err != nil {
			cherryLogger.Error("JSON marshal error data:%v", pushStruct)
			return
		}
		push.Data = jsonData
		pc.IncomingMsgChan <- push
	}
}

// NewClient returns a new client with the auto documentation route.
func NewClient() *Client {

	pc := newClient()
	// 设置服务器push过来消息的callback
	pc.OnUnexpectedEvent(UnexpectedEventCb(pc))
	return pc
}

// Connect to server
func (pc *Client) Connect(addr string) error {
	return pc.ConnectToTCP(addr)
}

// ConnectWS to websocket server
// func (pc *Client) ConnectWS(addr string) error {
// 	return pc.Connector.StartWS(addr)
// }

// SendRequest sends a request to the server
func (pc *Client) SendRequest(route string, data []byte) (uint, error) {
	requestStruct, err := routeMessage(route)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, requestStruct)
	if err != nil {
		return 0, err
	}

	response, err := pc.Request(route, requestStruct)

	responseStruct, err := routeMessage(response.Route)
	if err != nil {
		cherryLogger.Error(err.Error())
		return 0, err
	}
	err = pc.Serializer().Unmarshal(response.Data, responseStruct)
	if err != nil {
		cherryLogger.Error("unmarshal error data:%v", response.Data)
		return 0, err
	}
	jsonData, err := json.Marshal(responseStruct)
	if err != nil {
		cherryLogger.Error("JSON marshal error data:%v", responseStruct)
		return 0, err
	}
	response.Data = jsonData
	pc.IncomingMsgChan <- response
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// SendNotify sends a notify to the server
func (pc *Client) SendNotify(route string, data []byte) error {
	notifyStruct, err := routeMessage(route)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, notifyStruct)
	if err != nil {
		return err
	}
	return pc.Notify(route, notifyStruct)
}
