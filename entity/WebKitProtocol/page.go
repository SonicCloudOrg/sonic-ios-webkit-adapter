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

type Setting string

type ResourceType string

type CoordinateSystem string

type CookieSameSitePolicy string

type Appearance string

type Frame struct {
	Id             *string   `json:"id"`
	ParentId       *string   `json:"parentId,omitempty"`
	LoaderId       *LoaderId `json:"loaderId"`
	Name           *string   `json:"name,omitempty"`
	Url            *string   `json:"url"`
	SecurityOrigin *string   `json:"securityOrigin"`
	MimeType       *string   `json:"mimeType"`
}

type FrameResource struct {
	Url          *string       `json:"url"`
	Type         *ResourceType `json:"type"`
	MimeType     *string       `json:"mimeType"`
	Failed       *bool         `json:"failed,omitempty"`
	Canceled     *bool         `json:"canceled,omitempty"`
	SourceMapURL *string       `json:"sourceMapURL,omitempty"`
	TargetId     *string       `json:"targetId,omitempty"`
}

type FrameResourceTree struct {
	Frame       *Frame              `json:"frame"`
	ChildFrames []FrameResourceTree `json:"childFrames,omitempty"`
	Resources   []FrameResource     `json:"resources"`
}

type SearchResult struct {
	Url          *string    `json:"url"`
	FrameId      *FrameId   `json:"frameId"`
	MatchesCount *int       `json:"matchesCount"`
	RequestId    *RequestId `json:"requestId,omitempty"`
}

type Cookie struct {
	Name     *string               `json:"name"`
	Value    *string               `json:"value"`
	Domain   *string               `json:"domain"`
	Path     *string               `json:"path"`
	Expires  *int                  `json:"expires"`
	Session  *bool                 `json:"session"`
	HttpOnly *bool                 `json:"httpOnly"`
	Secure   *bool                 `json:"secure"`
	SameSite *CookieSameSitePolicy `json:"sameSite"`
}
