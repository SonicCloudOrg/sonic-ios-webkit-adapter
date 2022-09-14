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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
)

type iOS9 struct {
}

func initIOS9(protocol *protocolAdapter) {
	result := &iOS9{}
	protocol.init()
	protocol.mapSelectorList = result.mapSelectorList
}

func (i *iOS9) mapSelectorList(selectorList gjson.Result, message string) string {
	cssRange := selectorList.Get("range")
	var err error
	var newMsg = message
	var oldMsg = message
	for _, selector := range selectorList.Get("selectors").Array() {
		newMsg, err = sjson.Set(newMsg, selector.Path(oldMsg), map[string]interface{}{
			"text": selector.Value(),
		})
		if cssRange.Exists() {
			newMsg, err = sjson.Set(newMsg, selector.Get("range").Path(oldMsg), cssRange.Value())
		}
		if err != nil {
			log.Panic(err)
		}
	}
	newMsg, err = sjson.Delete(newMsg, cssRange.Path(oldMsg))
	if err != nil {
		log.Panic(err)
	}
	return newMsg
}
