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
package entity

type TargetProtocol struct {
	ID       int         `json:"id"`
	Method   string      `json:"method"`
	Params   interface{} `json:"params,omitempty"`
	TargetId string      `json:"targetId,omitempty"`
}

type TargetParams struct {
	ID       int           `json:"id,omitempty"`
	Message  interface{}   `json:"message,omitempty"`
	TargetId string        `json:"targetId,omitempty"`
	Edits    []interface{} `json:"edits,omitempty"`
}

type IRange struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type IDisabledStyle struct {
	Content  string
	CssRange IRange
}
