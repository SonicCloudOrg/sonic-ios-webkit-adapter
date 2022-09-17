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
	"github.com/SonicCloudOrg/sonic-ios-webkit-adapter/entity/WebKitProtocol"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/yezihack/e"
	"log"
)

type iOS8 struct {
	adapter *Adapter
}

func initIOS8(protocol *protocolAdapter) {
	result := &iOS8{
		adapter: protocol.adapter,
	}
	protocol.init()
	protocol.adapter.AddWebkitMessageFilter("error", result.targetError)
	protocol.mapSelectorList = result.mapSelectorList
}

func (i *iOS8) targetError(message []byte) []byte {
	params := map[string]interface{}{
		"id":     gjson.Get(string(message), "id"),
		"result": map[string]interface{}{},
	}
	msg, err := sjson.Set("", "", params)
	if err != nil {
		log.Println(e.Convert(err).ToStr())
	}
	return []byte(msg)
}

func (i *iOS8) mapSelectorList(selectorList *WebKitProtocol.SelectorList) {
	cssRange := selectorList.Range
	for index, _ := range selectorList.Selectors {
		selectorList.Selectors[index] = WebKitProtocol.CSSSelector{
			Text: selectorList.Selectors,
		}
		if cssRange != nil {
			selectorList.Selectors[index].Range = cssRange
		}
	}
	selectorList.Range = nil
}
