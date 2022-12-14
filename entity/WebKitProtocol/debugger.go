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

type BreakpointId string

type BreakpointActionIdentifier int

type ScriptId string

type CallFrameId string

type Location struct {
	ScriptId     *ScriptId `json:"scriptId"`
	LineNumber   *int      `json:"lineNumber"`
	ColumnNumber *int      `json:"columnNumber,omitempty"`
}

type BreakpointAction struct {
	Type               *string                     `json:"type"`
	Data               *string                     `json:"data,omitempty"`
	Id                 *BreakpointActionIdentifier `json:"id,omitempty"`
	EmulateUserGesture *bool                       `json:"emulateUserGesture,omitempty"`
}

type BreakpointOptions struct {
	Condition    *string            `json:"condition,omitempty"`
	Actions      []BreakpointAction `json:"actions,omitempty"`
	AutoContinue *bool              `json:"autoContinue,omitempty"`
	IgnoreCount  *int               `json:"ignoreCount,omitempty"`
}

type FunctionDetails struct {
	Location    *Location `json:"location"`
	Name        *string   `json:"name,omitempty"`
	DisplayName *string   `json:"displayName,omitempty"`
	ScopeChain  []Scope   `json:"scopeChain,omitempty"`
}

type DebuggerCallFrame struct {
	CallFrameId   *CallFrameId  `json:"callFrameId"`
	FunctionName  *string       `json:"functionName"`
	Location      *Location     `json:"location"`
	ScopeChain    []Scope       `json:"scopeChain"`
	This          *RemoteObject `json:"this"`
	IsTailDeleted *bool         `json:"isTailDeleted"`
}

type Scope struct {
	Object   *RemoteObject `json:"object"`
	Type     *string       `json:"type"`
	Name     *string       `json:"name,omitempty"`
	Location *Location     `json:"location,omitempty"`
	Empty    *bool         `json:"empty,omitempty"`
}

type ProbeSample struct {
	ProbeId   *BreakpointActionIdentifier `json:"probeId"`
	SampleId  *int                        `json:"sampleId"`
	BatchId   *int                        `json:"batchId"`
	Timestamp *int                        `json:"timestamp"`
	Payload   *RemoteObject               `json:"payload"`
}

type AssertPauseReason struct {
	Message *string `json:"message,omitempty"`
}

type BreakpointPauseReason struct {
	BreakpointId *string `json:"breakpointId"`
}

type CSPViolationPauseReason struct {
	Directive *string `json:"directive"`
}
