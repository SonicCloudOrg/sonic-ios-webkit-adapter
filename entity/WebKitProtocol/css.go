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
	PseudoId *PseudoId   `json:"pseudoId"`
	Matches  []RuleMatch `json:"matches"`
}

type InheritedStyleEntry struct {
	InlineStyle     *CSSStyle   `json:"inlineStyle,omitempty"`
	MatchedCSSRules []RuleMatch `json:"matchedCSSRules"`
}

type RuleMatch struct {
	Rule              *CSSRule `json:"rule"`
	MatchingSelectors []int    `json:"matchingSelectors"`
}

type CSSSelector struct {
	Text        interface{} `json:"text"`
	Specificity []int       `json:"specificity,omitempty"`
	Dynamic     *bool       `json:"dynamic,omitempty"`
	// devtool
	Range interface{} `json:"range,omitempty"`
}

type SelectorList struct {
	Selectors []CSSSelector `json:"selectors"`
	Text      *string       `json:"text,omitempty"`
	Range     *SourceRange  `json:"range,omitempty"`
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
	Rules        []CSSRule     `json:"rules"`
	Text         *string       `json:"text,omitempty"`
}

type CSSRule struct {
	RuleId       *CSSRuleId        `json:"ruleId,omitempty"`
	SelectorList *SelectorList     `json:"selectorList,omitempty"`
	SourceURL    *string           `json:"sourceURL,omitempty"`
	SourceLine   *int              `json:"sourceLine,omitempty"`
	Origin       *StyleSheetOrigin `json:"origin,omitempty"`
	Style        *CSSStyle         `json:"style,omitempty"`
	Groupings    []Grouping        `json:"groupings,omitempty"`

	// chrome devtool field
	DevToolStyleSheetId *StyleSheetId `json:"styleSheetId"`
}

type SourceRange struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
	// devtool
	Content *string      `json:"content,omitempty"`
	Range   *SourceRange `json:"range,omitempty"`
}

type ShorthandEntry struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type CSSPropertyInfo struct {
	Name      *string  `json:"name"`
	Aliases   []string `json:"aliases,omitempty"`
	Longhands []string `json:"longhands,omitempty"`
	Values    []string `json:"values,omitempty"`
	Inherited *bool    `json:"inherited,omitempty"`
}

type CSSComputedStyleProperty struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type CSSStyle struct {
	StyleId          *CSSStyleId      `json:"styleId,omitempty"`
	CssProperties    []CSSProperty    `json:"cssProperties"`
	ShorthandEntries []ShorthandEntry `json:"shorthandEntries"`
	CssText          *string          `json:"cssText,omitempty"`
	Range            *SourceRange     `json:"range,omitempty"`
	Width            *string          `json:"width,omitempty"`
	Height           *string          `json:"height,omitempty"`
	// devtool
	StyleSheetId *StyleSheetId `json:"styleSheetId,omitempty"`
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
	// devtool
	Disabled *bool `json:"disabled,omitempty"`
}

type Grouping struct {
	Type      *string `json:"type"`
	Text      *string `json:"text,omitempty"`
	SourceURL *string `json:"sourceURL,omitempty"`
}

type Font struct {
	DisplayName        *string             `json:"displayName"`
	VariationAxes      []FontVariationAxis `json:"variationAxes"`
	SynthesizedBold    *bool               `json:"synthesizedBold,omitempty"`
	SynthesizedOblique *bool               `json:"synthesizedOblique,omitempty"`
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

// result

type GetMatchedStylesForNodeResult struct {
	MatchedCSSRules []RuleMatch           `json:"matchedCSSRules,omitempty"`
	PseudoElements  []PseudoIdMatches     `json:"pseudoElements,omitempty"`
	Inherited       []InheritedStyleEntry `json:"inherited,omitempty"`
}

type GetInlineStylesForNodeResult struct {
	InlineStyle     *CSSStyle `json:"inlineStyle,omitempty"`
	AttributesStyle *CSSStyle `json:"attributesStyle,omitempty"`
}

type GetComputedStyleForNodeResult struct {
	ComputedStyle []CSSComputedStyleProperty `json:"computedStyle"`
}

type GetFontDataForNodeResult struct {
	PrimaryFont *Font `json:"primaryFont"`
}

type GetAllStyleSheetsResult struct {
	Headers []CSSStyleSheetHeader `json:"headers"`
}

type GetStyleSheetResult struct {
	StyleSheet *CSSStyleSheetBody `json:"styleSheet"`
}

type GetStyleSheetTextResult struct {
	Text *string `json:"text"`
}

type SetStyleTextResult struct {
	Style *CSSStyle `json:"style"`
}

type SetRuleSelectorResult struct {
	Rule *CSSRule `json:"rule"`
}

type CreateStyleSheetResult struct {
	StyleSheetId *StyleSheetId `json:"styleSheetId"`
}

type AddRuleResult struct {
	Rule *CSSRule `json:"rule"`
}

type GetSupportedCSSPropertiesResult struct {
	CssProperties []CSSPropertyInfo `json:"cssProperties"`
}

type GetSupportedSystemFontFamilyNamesResult struct {
	FontFamilyNames []string `json:"fontFamilyNames"`
}
