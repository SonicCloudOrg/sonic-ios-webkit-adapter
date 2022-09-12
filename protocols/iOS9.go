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
package protocols

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
)

type iOS9 struct {
}

func initIOS9(protocol *ProtocolAdapter) {
	result := &iOS9{}
	protocol.init()
	protocol.mapSelectorList = result.mapSelectorList
}

func (i *iOS9) mapSelectorList(selectorList gjson.Result, message string) string {
	cssRange := selectorList.Get("range")
	var err error
	for _, selector := range selectorList.Get("selectors").Array() {
		message, err = sjson.Set(message, selector.Path(message), map[string]interface{}{
			"text": selector.Value(),
		})
		if cssRange.Exists() {
			message, err = sjson.Set(message, selector.Get("range").Path(message), cssRange.Value())
		}
		if err != nil {
			log.Panic(err)
		}
	}
	message, err = sjson.Delete(message, cssRange.Path(message))
	if err != nil {
		log.Panic(err)
	}
	return message
}
