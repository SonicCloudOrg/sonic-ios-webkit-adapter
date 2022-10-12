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

type NodeId int

type EventListenerId int

type PseudoType string

type ShadowRootType string

type CustomElementState string

type LiveRegionRelevant string

type Node struct {
	NodeId                    *NodeId             `json:"nodeId"`
	NodeType                  *int                `json:"nodeType"`
	NodeName                  *string             `json:"nodeName"`
	LocalName                 *string             `json:"localName"`
	NodeValue                 *string             `json:"nodeValue"`
	FrameId                   *FrameId            `json:"frameId,omitempty"`
	ChildNodeCount            *int                `json:"childNodeCount,omitempty"`
	Children                  []Node              `json:"children,omitempty"`
	Attributes                []string            `json:"attributes,omitempty"`
	DocumentURL               *string             `json:"documentURL,omitempty"`
	BaseURL                   *string             `json:"baseURL,omitempty"`
	PublicId                  *string             `json:"publicId,omitempty"`
	SystemId                  *string             `json:"systemId,omitempty"`
	XmlVersion                *string             `json:"xmlVersion,omitempty"`
	Name                      *string             `json:"name,omitempty"`
	Value                     *string             `json:"value,omitempty"`
	PseudoType                *PseudoType         `json:"pseudoType,omitempty"`
	ShadowRootType            *ShadowRootType     `json:"shadowRootType,omitempty"`
	CustomElementState        *CustomElementState `json:"customElementState,omitempty"`
	ContentDocument           *Node               `json:"contentDocument,omitempty"`
	ShadowRoots               []Node              `json:"shadowRoots,omitempty"`
	TemplateContent           *Node               `json:"templateContent,omitempty"`
	PseudoElements            []Node              `json:"pseudoElements,omitempty"`
	ContentSecurityPolicyHash *string             `json:"contentSecurityPolicyHash,omitempty"`
	LayoutFlags               []string            `json:"layoutFlags,omitempty"`
}

type DataBinding struct {
	Binding *string `json:"binding"`
	Type    *string `json:"type,omitempty"`
	Value   *string `json:"value"`
}

type EventListener struct {
	EventListenerId *EventListenerId `json:"eventListenerId"`
	Type            *string          `json:"type"`
	UseCapture      *bool            `json:"useCapture"`
	IsAttribute     *bool            `json:"isAttribute"`
	NodeId          *NodeId          `json:"nodeId,omitempty"`
	OnWindow        *bool            `json:"onWindow,omitempty"`
	Location        *Location        `json:"location,omitempty"`
	HandlerName     *string          `json:"handlerName,omitempty"`
	Passive         *bool            `json:"passive,omitempty"`
	Once            *bool            `json:"once,omitempty"`
	Disabled        *bool            `json:"disabled,omitempty"`
	HasBreakpoint   *bool            `json:"hasBreakpoint,omitempty"`
}

type AccessibilityProperties struct {
	ActiveDescendantNodeId *NodeId  `json:"activeDescendantNodeId,omitempty"`
	Busy                   *bool    `json:"busy,omitempty"`
	Checked                *string  `json:"checked,omitempty"`
	ChildNodeIds           []NodeId `json:"childNodeIds,omitempty"`
	ControlledNodeIds      []NodeId `json:"controlledNodeIds,omitempty"`
	Current                *string  `json:"current,omitempty"`
	Disabled               *bool    `json:"disabled,omitempty"`
	HeadingLevel           *int     `json:"headingLevel,omitempty"`
	HierarchyLevel         *int     `json:"hierarchyLevel,omitempty"`
	IsPopUpButton          *bool    `json:"isPopUpButton,omitempty"`
	Exists                 *bool    `json:"exists"`
	Expanded               *bool    `json:"expanded,omitempty"`
	FlowedNodeIds          []NodeId `json:"flowedNodeIds,omitempty"`
	Focused                *bool    `json:"focused,omitempty"`
	Ignored                *bool    `json:"ignored,omitempty"`
	IgnoredByDefault       *bool    `json:"ignoredByDefault,omitempty"`
	Invalid                *string  `json:"invalid,omitempty"`
	Hidden                 *bool    `json:"hidden,omitempty"`
	Label                  *string  `json:"label"`
	LiveRegionAtomic       *bool    `json:"liveRegionAtomic,omitempty"`
	LiveRegionRelevant     []string `json:"liveRegionRelevant,omitempty"`
	LiveRegionStatus       *string  `json:"liveRegionStatus,omitempty"`
	MouseEventNodeId       *NodeId  `json:"mouseEventNodeId,omitempty"`
	NodeId                 *NodeId  `json:"nodeId"`
	OwnedNodeIds           []NodeId `json:"ownedNodeIds,omitempty"`
	ParentNodeId           *NodeId  `json:"parentNodeId,omitempty"`
	Pressed                *bool    `json:"pressed,omitempty"`
	Readonly               *bool    `json:"readonly,omitempty"`
	Required               *bool    `json:"required,omitempty"`
	Role                   *string  `json:"role"`
	Selected               *bool    `json:"selected,omitempty"`
	SelectedChildNodeIds   []NodeId `json:"selectedChildNodeIds,omitempty"`
}

type RGBAColor struct {
	R *int `json:"r"`
	G *int `json:"g"`
	B *int `json:"b"`
	A *int `json:"a,omitempty"`
}

type Quad []int

type HighlightConfig struct {
	ShowInfo     *bool      `json:"showInfo,omitempty"`
	ContentColor *RGBAColor `json:"contentColor,omitempty"`
	PaddingColor *RGBAColor `json:"paddingColor,omitempty"`
	BorderColor  *RGBAColor `json:"borderColor,omitempty"`
	MarginColor  *RGBAColor `json:"marginColor,omitempty"`
}

type Styleable struct {
	NodeId   *NodeId   `json:"nodeId"`
	PseudoId *PseudoId `json:"pseudoId,omitempty"`
}

// result
type GetDocumentResult struct {
	Root *Node `json:"root"`
}

type QuerySelectorResult struct {
	NodeId *NodeId `json:"nodeId,omitempty"`
}

type QuerySelectorAllResult struct {
	NodeIds []NodeId `json:"nodeIds"`
}

type SetNodeNameResult struct {
	NodeId *NodeId `json:"nodeId"`
}

type GetSupportedEventNamesResult struct {
	EventNames []string `json:"eventNames"`
}

type GetDataBindingsForNodeResult struct {
	DataBindings []DataBinding `json:"dataBindings"`
}

type GetAssociatedDataForNodeResult struct {
	AssociatedData *string `json:"associatedData,omitempty"`
}

type GetEventListenersForNodeResult struct {
	Listeners []EventListener `json:"listeners"`
}

type GetAccessibilityPropertiesForNodeResult struct {
	Properties *AccessibilityProperties `json:"properties"`
}

type GetOuterHTMLResult struct {
	OuterHTML *string `json:"outerHTML"`
}

type PerformSearchResult struct {
	SearchId    *string `json:"searchId"`
	ResultCount *int    `json:"resultCount"`
}

type GetSearchResultsResult struct {
	NodeIds []NodeId `json:"nodeIds"`
}

type RequestNodeResult struct {
	NodeId *NodeId `json:"nodeId"`
}

type PushNodeByPathToFrontendResult struct {
	NodeId *NodeId `json:"nodeId"`
}

type ResolveNodeResult struct {
	Object *RemoteObject `json:"object"`
}

type GetAttributesResult struct {
	Attributes []string `json:"attributes"`
}

type MoveToResult struct {
	NodeId *NodeId `json:"nodeId"`
}
