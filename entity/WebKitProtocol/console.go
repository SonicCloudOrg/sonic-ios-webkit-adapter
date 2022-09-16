package WebKitProtocol

type ChannelSource string

type ChannelLevel string

type Channel struct {
	Source *ChannelSource `json:"source"`
	Level  *ChannelLevel  `json:"level"`
}

type ConsoleMessage struct {
	Source           *ChannelSource  `json:"source"`
	Level            *string         `json:"level"`
	Text             *string         `json:"text"`
	Type             *string         `json:"type,omitempty"`
	Url              *string         `json:"url,omitempty"`
	Line             *int            `json:"line,omitempty"`
	Column           *int            `json:"column,omitempty"`
	RepeatCount      *int            `json:"repeatCount,omitempty"`
	Parameters       *[]RemoteObject `json:"parameters,omitempty"`
	StackTrace       *StackTrace     `json:"stackTrace,omitempty"`
	NetworkRequestId *RequestId      `json:"networkRequestId,omitempty"`
	Timestamp        *int            `json:"timestamp,omitempty"`
}

type ConsoleCallFrame struct {
	FunctionName *string   `json:"functionName"`
	Url          *string   `json:"url"`
	ScriptId     *ScriptId `json:"scriptId"`
	LineNumber   *int      `json:"lineNumber"`
	ColumnNumber *int      `json:"columnNumber"`
}

type StackTrace struct {
	CallFrames             *[]ConsoleCallFrame `json:"callFrames"`
	TopCallFrameIsBoundary *bool               `json:"topCallFrameIsBoundary,omitempty"`
	Truncated              *bool               `json:"truncated,omitempty"`
	ParentStackTrace       *StackTrace         `json:"parentStackTrace,omitempty"`
}
