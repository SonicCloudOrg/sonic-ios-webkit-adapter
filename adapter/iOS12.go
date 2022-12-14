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
)

type iOS12 struct {
	adapter *Adapter
}

func initIOS12(protocol *protocolAdapter) {
	protocol.adapter.SetTargetBased(true)
	result := &iOS12{
		adapter: protocol.adapter,
	}
	//protocol.init()
	protocol.adapter.AddWebkitMessageFilter("Target.targetCreated", result.targetCreated)
	initIOS9(protocol)
}

func (i *iOS12) targetCreated(message []byte) []byte {
	i.adapter.SetTargetID(gjson.Get(string(message), "params.targetInfo.targetId").String())
	return message
}
