package WebKitProtocol

type StyleSheetId string

type CSSStyleId struct {
	StyleSheetId *StyleSheetId `json:"styleSheetId"`
	Ordinal      *int          `json:"ordinal"`
}

type StyleSheetOrigin string

type CSSRuleId struct {
	StyleSheetId *StyleSheetId `json:"styleSheetId"`
	Ordinal      *int          `json:"ordinal"`
}

type PseudoId string

type ForceablePseudoClass string

type PseudoIdMatches struct {
	PseudoId *PseudoId    `json:"pseudoId"`
	Matches  *[]RuleMatch `json:"matches"`
}

type InheritedStyleEntry struct {
	InlineStyle     *CSSStyle    `json:"inlineStyle,omitempty"`
	MatchedCSSRules *[]RuleMatch `json:"matchedCSSRules"`
}

type RuleMatch struct {
	Rule              *CSSRule `json:"rule"`
	MatchingSelectors *[]int   `json:"matchingSelectors"`
}

type CSSSelector struct {
	Text        *string `json:"text"`
	Specificity *[]int  `json:"specificity,omitempty"`
	Dynamic     *bool   `json:"dynamic,omitempty"`
}

type SelectorList struct {
	Selectors *[]CSSSelector `json:"selectors"`
	Text      *string        `json:"text"`
	Range     *SourceRange   `json:"range,omitempty"`
}

type CSSStyleAttribute struct {
	Name  *string   `json:"name"`
	Style *CSSStyle `json:"style"`
}

type CSSStyleSheetHeader struct {
	StyleSheetId *StyleSheetId     `json:"styleSheetId"`
	FrameId      *FrameId          `json:"frameId"`
	SourceURL    *string           `json:"sourceURL"`
	Origin       *StyleSheetOrigin `json:"origin"`
	Title        *string           `json:"title"`
	Disabled     *bool             `json:"disabled"`
	IsInline     *bool             `json:"isInline"`
	StartLine    *int              `json:"startLine"`
	StartColumn  *int              `json:"startColumn"`
}

type CSSStyleSheetBody struct {
	StyleSheetId *StyleSheetId `json:"styleSheetId"`
	Rules        *[]CSSRule    `json:"rules"`
	Text         *string       `json:"text,omitempty"`
}

type CSSRule struct {
	RuleId       *CSSRuleId        `json:"ruleId,omitempty"`
	SelectorList *SelectorList     `json:"selectorList"`
	SourceURL    *string           `json:"sourceURL,omitempty"`
	SourceLine   *int              `json:"sourceLine"`
	Origin       *StyleSheetOrigin `json:"origin"`
	Style        *CSSStyle         `json:"style"`
	Groupings    *[]Grouping       `json:"groupings,omitempty"`
}

type SourceRange struct {
	StartLine   *int `json:"startLine"`
	StartColumn *int `json:"startColumn"`
	EndLine     *int `json:"endLine"`
	EndColumn   *int `json:"endColumn"`
}

type ShorthandEntry struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type CSSPropertyInfo struct {
	Name      *string   `json:"name"`
	Aliases   *[]string `json:"aliases,omitempty"`
	Longhands *[]string `json:"longhands,omitempty"`
	Values    *[]string `json:"values,omitempty"`
	Inherited *bool     `json:"inherited,omitempty"`
}

type CSSComputedStyleProperty struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type CSSStyle struct {
	StyleId          *CSSStyleId       `json:"styleId,omitempty"`
	CssProperties    *[]CSSProperty    `json:"cssProperties"`
	ShorthandEntries *[]ShorthandEntry `json:"shorthandEntries"`
	CssText          *string           `json:"cssText,omitempty"`
	Range            *SourceRange      `json:"range,omitempty"`
	Width            *string           `json:"width,omitempty"`
	Height           *string           `json:"height,omitempty"`
}

type CSSPropertyStatus string

type CSSProperty struct {
	Name     *string            `json:"name"`
	Value    *string            `json:"value"`
	Priority *string            `json:"priority,omitempty"`
	Implicit *bool              `json:"implicit,omitempty"`
	Text     *string            `json:"text,omitempty"`
	ParsedOk *bool              `json:"parsedOk,omitempty"`
	Status   *CSSPropertyStatus `json:"status,omitempty"`
	Range    *SourceRange       `json:"range,omitempty"`
}

type Grouping struct {
	Type      *string `json:"type"`
	Text      *string `json:"text,omitempty"`
	SourceURL *string `json:"sourceURL,omitempty"`
}

type Font struct {
	DisplayName        *string              `json:"displayName"`
	VariationAxes      *[]FontVariationAxis `json:"variationAxes"`
	SynthesizedBold    *bool                `json:"synthesizedBold,omitempty"`
	SynthesizedOblique *bool                `json:"synthesizedOblique,omitempty"`
}

type FontVariationAxis struct {
	Name         *string `json:"name,omitempty"`
	Tag          *string `json:"tag"`
	MinimumValue *int    `json:"minimumValue"`
	MaximumValue *int    `json:"maximumValue"`
	DefaultValue *int    `json:"defaultValue"`
}

type LayoutFlag string

type LayoutContextTypeChangedMode string
