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

package repl2

import (
	"encoding/json"
	"errors"
	"strings"

	cherryLogger "github.com/cherry-game/cherry/logger"
	cherryClient "github.com/cherry-game/cherry/net/client"
	cmessage "github.com/cherry-game/cherry/net/message"
	cherrySerializer "github.com/cherry-game/cherry/net/serializer"
)

func connect(logger Log, addr string, onMessageCallback func([]byte)) (err error) {
	if pClient != nil && pClient.IsConnected() {
		return errors.New("already connected")
	}

	switch {
	case docsString != "":
		//err = protoClient(logger, addr)
	default:
		logger.Println("Using json client")
		pClient = cherryClient.New(cherryClient.WithSerializer(cherrySerializer.NewJSON()))
	}
	pClient.OnUnexpectedEvent(UnexpectedEventCb(pClient))

	if err != nil {
		return err
	}

	if err = tryConnect(addr); err != nil {
		logger.Println("Failed to connect!")
		return err
	}

	logger.Println("connected!")
	disconnectedCh = make(chan bool, 1)

	return nil
}

// UnexpectedEventCb returns a function to deal with un listened event
func UnexpectedEventCb(pc *cherryClient.Client) func(data interface{}) {
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
		logger.Println(string(jsonData))
	}
}
func push(logger Log, args []string) error {
	if pClient != nil {
		return errors.New("use this command before connect")
	}

	if len(args) != 2 {
		return errors.New(`push should be in the format: push {route} {type}`)
	}

	route := args[0]
	pushType := args[1]

	if docsString == "" {
		logger.Println("Only for probuffer servers")
		return nil
	}

	pushInfo[route] = pushType

	return nil
}

func request(args []string) error {
	if pClient == nil {
		return errors.New("not connected")
	}

	if !pClient.IsConnected() {
		return errors.New("not connected")
	}

	if len(args) < 1 {
		return errors.New(`request should be in the format: request {route} [data]`)
	}

	route := args[0]

	var data []byte
	if len(args) > 1 {
		data = []byte(strings.Join(args[1:], ""))
	}

	response, err := pClient.Request(route, data)
	if err != nil {
		return err
	}

	responseStruct, err := routeMessage(response.Route)
	if err != nil {
		cherryLogger.Error(err.Error())
		return err
	}
	err = pClient.Serializer().Unmarshal(response.Data, responseStruct)
	if err != nil {
		cherryLogger.Error("unmarshal error data:%v", response.Data)
		return err
	}
	jsonData, err := json.Marshal(responseStruct)
	if err != nil {
		cherryLogger.Error("JSON marshal error data:%v", responseStruct)
		return err
	}
	logger.Println(string(jsonData))
	return nil
}

func notify(logger Log, args []string) error {
	if pClient == nil {
		return errors.New("not connected")
	}

	if !pClient.IsConnected() {
		return errors.New("not connected")
	}

	if len(args) < 1 {
		return errors.New(`notify should be in the format: notify {route} [data]`)
	}

	route := args[0]
	var data []byte
	if len(args) > 1 {
		data = []byte(strings.Join(args[1:], ""))
	}

	if err := pClient.Notify(route, data); err != nil {
		return err
	}

	return nil
}

func disconnect() {
	if pClient.IsConnected() {
		disconnectedCh <- true
		pClient.Disconnect()
	}
}
