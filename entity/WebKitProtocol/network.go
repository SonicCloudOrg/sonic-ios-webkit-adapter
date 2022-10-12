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

type LoaderId string

type FrameId string

type RequestId string

type Timestamp int

type Walltime int

type ReferrerPolicy string

type Headers interface{}

type ResourceTiming struct {
	StartTime             *Timestamp `json:"startTime"`
	RedirectStart         *Timestamp `json:"redirectStart"`
	RedirectEnd           *Timestamp `json:"redirectEnd"`
	FetchStart            *Timestamp `json:"fetchStart"`
	DomainLookupStart     *int       `json:"domainLookupStart"`
	DomainLookupEnd       *int       `json:"domainLookupEnd"`
	ConnectStart          *int       `json:"connectStart"`
	ConnectEnd            *int       `json:"connectEnd"`
	SecureConnectionStart *int       `json:"secureConnectionStart"`
	RequestStart          *int       `json:"requestStart"`
	ResponseStart         *int       `json:"responseStart"`
	ResponseEnd           *int       `json:"responseEnd"`
}

type Request struct {
	Url            *string         `json:"url"`
	Method         *string         `json:"method"`
	Headers        *Headers        `json:"headers"`
	PostData       *string         `json:"postData,omitempty"`
	ReferrerPolicy *ReferrerPolicy `json:"referrerPolicy,omitempty"`
	Integrity      *string         `json:"integrity,omitempty"`
}

type Response struct {
	Url            *string         `json:"url"`
	Status         *int            `json:"status"`
	StatusText     *string         `json:"statusText"`
	Headers        *Headers        `json:"headers"`
	MimeType       *string         `json:"mimeType"`
	Source         *string         `json:"source"`
	RequestHeaders *Headers        `json:"requestHeaders,omitempty"`
	Timing         *ResourceTiming `json:"timing,omitempty"`
	Security       *Security       `json:"security,omitempty"`
}

type Metrics struct {
	Protocol                    *string     `json:"protocol,omitempty"`
	Priority                    *string     `json:"priority,omitempty"`
	ConnectionIdentifier        *string     `json:"connectionIdentifier,omitempty"`
	RemoteAddress               *string     `json:"remoteAddress,omitempty"`
	RequestHeaders              *Headers    `json:"requestHeaders,omitempty"`
	RequestHeaderBytesSent      *int        `json:"requestHeaderBytesSent,omitempty"`
	RequestBodyBytesSent        *int        `json:"requestBodyBytesSent,omitempty"`
	ResponseHeaderBytesReceived *int        `json:"responseHeaderBytesReceived,omitempty"`
	ResponseBodyBytesReceived   *int        `json:"responseBodyBytesReceived,omitempty"`
	ResponseBodyDecodedSize     *int        `json:"responseBodyDecodedSize,omitempty"`
	SecurityConnection          *Connection `json:"securityConnection,omitempty"`
	IsProxyConnection           *bool       `json:"isProxyConnection,omitempty"`
}

type WebSocketRequest struct {
	Headers *Headers `json:"headers"`
}

type WebSocketResponse struct {
	Status     *int     `json:"status"`
	StatusText *string  `json:"statusText"`
	Headers    *Headers `json:"headers"`
}

type WebSocketFrame struct {
	Opcode        *int    `json:"opcode"`
	Mask          *bool   `json:"mask"`
	PayloadData   *string `json:"payloadData"`
	PayloadLength *int    `json:"payloadLength"`
}

type CachedResource struct {
	Url          *string       `json:"url"`
	Type         *ResourceType `json:"type"`
	Response     *Response     `json:"response,omitempty"`
	BodySize     *int          `json:"bodySize"`
	SourceMapURL *string       `json:"sourceMapURL,omitempty"`
}

type Initiator struct {
	Type       *string     `json:"type"`
	StackTrace *StackTrace `json:"stackTrace,omitempty"`
	Url        *string     `json:"url,omitempty"`
	LineNumber *int        `json:"lineNumber,omitempty"`
	NodeId     *NodeId     `json:"nodeId,omitempty"`
}

type NetworkStage string

type ResourceErrorType string
