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
package WebKitProtocol

type ChannelSource string

type ChannelLevel string

type Channel struct {
	Source *ChannelSource `json:"source"`
	Level  *ChannelLevel  `json:"level"`
}

type ConsoleMessage struct {
	Source           *ChannelSource `json:"source"`
	Level            *string        `json:"level"`
	Text             *string        `json:"text"`
	Type             *string        `json:"type,omitempty"`
	Url              *string        `json:"url,omitempty"`
	Line             *int           `json:"line,omitempty"`
	Column           *int           `json:"column,omitempty"`
	RepeatCount      *int           `json:"repeatCount,omitempty"`
	Parameters       []RemoteObject `json:"parameters,omitempty"`
	StackTrace       *StackTrace    `json:"stackTrace,omitempty"`
	NetworkRequestId *RequestId     `json:"networkRequestId,omitempty"`
	Timestamp        *int           `json:"timestamp,omitempty"`
}

type ConsoleCallFrame struct {
	FunctionName *string   `json:"functionName"`
	Url          *string   `json:"url"`
	ScriptId     *ScriptId `json:"scriptId"`
	LineNumber   *int      `json:"lineNumber"`
	ColumnNumber *int      `json:"columnNumber"`
}

type StackTrace struct {
	CallFrames             []ConsoleCallFrame `json:"callFrames"`
	TopCallFrameIsBoundary *bool              `json:"topCallFrameIsBoundary,omitempty"`
	Truncated              *bool              `json:"truncated,omitempty"`
	ParentStackTrace       *StackTrace        `json:"parentStackTrace,omitempty"`
}
