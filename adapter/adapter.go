/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package adapters

import (
	"encoding/json"
	"github.com/SonicCloudOrg/sonic-ios-webkit-adapter/entity"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)

type MessageAdapters func(message []byte) []byte

type Adapter struct {
	targetID          string
	messageFilters    map[string]MessageAdapters
	messageBuffer     [][]byte
	isTargetBased     bool
	applicationID     *string
	pageID            *int
	waitingForID      int
	toolRequestMap    map[int64]string
	adapterRequestMap map[int64]func(message []byte)
	wsToolServer      *websocket.Conn
	wsWebkitServer    *websocket.Conn
	isToolConnect     bool
	// 给iOS
	sendWebkit func([]byte)
	// 给devtool
	sendDevTool func([]byte)
	// recv for IOS
	receiveWebKit func([]byte)
	// recv for devtool
	receiveDevTool func([]byte)
}

func NewAdapter(wsToolServer *websocket.Conn, version string) *Adapter {
	adapter := &Adapter{
		messageFilters:    make(map[string]MessageAdapters),
		toolRequestMap:    make(map[int64]string),
		adapterRequestMap: make(map[int64]func(message []byte)),
	}
	adapter.sendWebkit = adapter.defaultSendWebkit
	adapter.receiveWebKit = adapter.defaultReceiveWebkit
	adapter.sendDevTool = adapter.defaultSendDevTool
	adapter.receiveDevTool = adapter.defaultReceiveDevTool

	initProtocolAdapter(adapter, version)

	return adapter
}

func (a *Adapter) addMessageFilter(method string, filter MessageAdapters) {
	if a.messageFilters == nil {
		a.messageFilters = make(map[string]MessageAdapters)
	}
	a.messageFilters[method] = filter
}

func (a *Adapter) CallTarget(method string, params interface{}, callFunc func(message []byte)) {
	a.waitingForID -= 1
	var message = &entity.TargetProtocol{}
	arr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	println(string(arr))
	message.ID = a.waitingForID
	message.Method = method
	message.Params = params
	if callFunc != nil {
		a.adapterRequestMap[int64(a.waitingForID)] = callFunc
	}
	a.sendToTarget(message)
}

func (a *Adapter) sendToTarget(message *entity.TargetProtocol) {
	log.Println("origin send message:")
	arr, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(arr))
	if a.isTargetBased {
		if !strings.Contains(message.Method, "Target") {
			var newMessage = &entity.TargetProtocol{}

			newMessage.ID = message.ID
			newMessage.Method = "Target.sendMessageToTarget"
			newMessage.Params = &entity.TargetParams{
				TargetId: a.targetID,
				Message:  string(arr),
				ID:       message.ID,
			}
			message = newMessage
			log.Println("new send message:")
		}
	}
	arr, err = json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	a.sendWebkit(arr)
}

func (a *Adapter) FireEventToTools(method string, params interface{}) {
	response := map[string]interface{}{
		"method": method,
		"params": params,
	}
	arr, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}
	a.sendDevTool(arr)
}

func (a *Adapter) FireResultToTools(id int, params interface{}) {
	response := map[string]interface{}{
		"id":     id,
		"result": params,
	}
	arr, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}
	a.sendDevTool(arr)
}

func (a *Adapter) ReplyWithEmpty(msg string) []byte {
	a.FireResultToTools(int(gjson.Get(msg, "id").Int()), map[string]interface{}{})
	return nil
}

func (a *Adapter) SetTargetBased(flag bool) {
	a.isTargetBased = flag
}

func (a *Adapter) SetTargetID(targetID string) {
	a.targetID = targetID
}

func (a *Adapter) SetIsConnect(flag bool) {
	a.isToolConnect = flag
}

// todo webkit debug ws close case
func (a *Adapter) Connect(wsPath string, toolWs *websocket.Conn) {
	a.wsToolServer = toolWs
	conn, _, err := websocket.DefaultDialer.Dial(wsPath, nil)
	if err != nil {
		log.Panic(err)
	}
	a.wsWebkitServer = conn
	a.SetIsConnect(true)

	go func() {
		for {
			_, message, err := a.wsWebkitServer.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}
			if message != nil {
				if len(message) == 0 {
					continue
				}
				a.receiveWebKit(message)
			}
		}
	}()
	for _, value := range a.messageBuffer {
		a.receiveDevTool(value)
	}
	a.messageBuffer = [][]byte{}
}

func (a *Adapter) SendMessageWebkit(message []byte) {
	a.sendWebkit(message)
}

func (a *Adapter) ReceiveMessageWebkit(message []byte) {
	a.receiveWebKit(message)
}

func (a *Adapter) SendMessageDevTool(message []byte) {
	a.sendDevTool(message)
}

func (a *Adapter) ReceiveMessageDevTool(message []byte) {
	a.receiveDevTool(message)
}

func (a *Adapter) defaultSendWebkit(message []byte) {
	if message == nil {
		return
	}
	err := a.wsWebkitServer.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Panic(err)
	}
}

func (a *Adapter) defaultReceiveWebkit(message []byte) {
	msg := string(message)
	if a.isTargetBased {
		method := gjson.Get(msg, "method")
		if !method.Exists() || !strings.Contains(method.String(), "Target") {
			return
		}
		if method.String() == "Target.dispatchMessageFromTarget" {
			msg = gjson.Get(msg, "params.message").String()
		}
	}
	// id exists in the message
	if gjson.Get(msg, "id").Exists() {
		id := gjson.Get(msg, "id").Int()
		if a.toolRequestMap[id] != "" {
			var eventName = a.toolRequestMap[id]
			if strings.Contains(msg, "err") && a.messageFilters["error"] != nil {
				eventName = "error"
			}

			if a.messageFilters[eventName] != nil {
				rawMessage := a.messageFilters[eventName]([]byte(msg))
				if rawMessage != nil {
					a.sendDevTool(rawMessage)
				}
			} else {
				a.sendDevTool([]byte(msg))
			}
		} else if a.adapterRequestMap[id] != nil {
			adapterFunc := a.adapterRequestMap[id]
			delete(a.adapterRequestMap, id)
			// 调用注册的回调函数
			if strings.Contains(msg, "result") {
				adapterFunc([]byte(gjson.Get(msg, "result").String()))
			} else if strings.Contains(msg, "error") {
				adapterFunc([]byte(gjson.Get(msg, "error").String()))
			} else {
				log.Println("unhandled type of request message from target:")
				log.Println(msg)
				log.Println()
			}
		} else {
			log.Println("unhandled message from target:")
			log.Println(msg)
			log.Println()
		}
	} else {
		var eventName = gjson.Get(msg, "method").String()
		if a.messageFilters[eventName] != nil {
			rawMessage := a.messageFilters[eventName]([]byte(msg))
			if rawMessage != nil {
				a.sendDevTool(rawMessage)
			}
		} else {
			a.sendDevTool([]byte(msg))
		}
	}
}

func (a *Adapter) defaultSendDevTool(message []byte) {
	if message == nil {
		return
	}
	err := a.wsToolServer.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}
}

func (a *Adapter) defaultReceiveDevTool(message []byte) {
	if !a.isToolConnect {
		a.messageBuffer = append(a.messageBuffer, message)
		return
	}
	msg := string(message)
	eventName := gjson.Get(msg, "method").String()
	id := gjson.Get(msg, "id").Int()
	a.toolRequestMap[id] = eventName

	if a.messageFilters[eventName] != nil {
		message = a.messageFilters[eventName](message)
	}
	if message != nil {
		protocolMessage := &entity.TargetProtocol{}
		err := json.Unmarshal(message, protocolMessage)
		if err != nil {
			log.Panic(err)
		}
		a.sendToTarget(protocolMessage)
	}
}

func (a *Adapter) SetSendWebkit(sendWebkit func([]byte)) {
	a.sendWebkit = sendWebkit
}

func (a *Adapter) SetSendDevTool(sendDevTool func([]byte)) {
	a.sendDevTool = sendDevTool
}

func (a *Adapter) SetReceiveWebkit(receiveWebkit func([]byte)) {
	a.receiveWebKit = receiveWebkit
}

func (a *Adapter) SetReceiveDevTool(receiveDevTool func([]byte)) {
	a.receiveDevTool = receiveDevTool
}
