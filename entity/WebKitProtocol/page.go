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
