package WebKitProtocol

type RemoteObjectId string

type RemoteObject struct {
	Type           *string         `json:"type"`
	Subtype        *string         `json:"subtype,omitempty"`
	ClassName      *string         `json:"className,omitempty"`
	Value          *any            `json:"value,omitempty"`
	Description    *string         `json:"description,omitempty"`
	ObjectId       *RemoteObjectId `json:"objectId,omitempty"`
	Size           *int            `json:"size,omitempty"`
	ClassPrototype *RemoteObject   `json:"classPrototype,omitempty"`
	Preview        *ObjectPreview  `json:"preview,omitempty"`
}

type ObjectPreview struct {
	Type        *string            `json:"type"`
	Subtype     *string            `json:"subtype,omitempty"`
	Description *string            `json:"description,omitempty"`
	Lossless    *bool              `json:"lossless"`
	Overflow    *bool              `json:"overflow,omitempty"`
	Properties  *[]PropertyPreview `json:"properties,omitempty"`
	Entries     *[]EntryPreview    `json:"entries,omitempty"`
	Size        *int               `json:"size,omitempty"`
}

type PropertyPreview struct {
	Name         *string        `json:"name"`
	Type         *string        `json:"type"`
	Subtype      *string        `json:"subtype,omitempty"`
	Value        *string        `json:"value,omitempty"`
	ValuePreview *ObjectPreview `json:"valuePreview,omitempty"`
	Internal     *bool          `json:"internal,omitempty"`
}

type EntryPreview struct {
	Key   *ObjectPreview `json:"key,omitempty"`
	Value *ObjectPreview `json:"value"`
}

type CollectionEntry struct {
	Key   *RemoteObject `json:"key,omitempty"`
	Value *RemoteObject `json:"value"`
}

type PropertyDescriptor struct {
	Name         *string       `json:"name"`
	Value        *RemoteObject `json:"value,omitempty"`
	Writable     *bool         `json:"writable,omitempty"`
	Get          *RemoteObject `json:"get,omitempty"`
	Set          *RemoteObject `json:"set,omitempty"`
	WasThrown    *bool         `json:"wasThrown,omitempty"`
	Configurable *bool         `json:"configurable,omitempty"`
	Enumerable   *bool         `json:"enumerable,omitempty"`
	IsOwn        *bool         `json:"isOwn,omitempty"`
	Symbol       *RemoteObject `json:"symbol,omitempty"`
	NativeGetter *bool         `json:"nativeGetter,omitempty"`
}

type InternalPropertyDescriptor struct {
	Name  *string       `json:"name"`
	Value *RemoteObject `json:"value,omitempty"`
}

type CallArgument struct {
	Value    *any            `json:"value,omitempty"`
	ObjectId *RemoteObjectId `json:"objectId,omitempty"`
}

type ExecutionContextId int

type ExecutionContextType string

type ExecutionContextDescription struct {
	Id      *ExecutionContextId   `json:"id"`
	Type    *ExecutionContextType `json:"type"`
	Name    *string               `json:"name"`
	FrameId *FrameId              `json:"frameId"`
}

type SyntaxErrorType string

type ErrorRange struct {
	StartOffset *int `json:"startOffset"`
	EndOffset   *int `json:"endOffset"`
}

type StructureDescription struct {
	Fields             *[]string             `json:"fields,omitempty"`
	OptionalFields     *[]string             `json:"optionalFields,omitempty"`
	ConstructorName    *string               `json:"constructorName,omitempty"`
	PrototypeStructure *StructureDescription `json:"prototypeStructure,omitempty"`
	IsImprecise        *bool                 `json:"isImprecise,omitempty"`
}

type TypeSet struct {
	IsFunction  *bool `json:"isFunction"`
	IsUndefined *bool `json:"isUndefined"`
	IsNull      *bool `json:"isNull"`
	IsBoolean   *bool `json:"isBoolean"`
	IsInteger   *bool `json:"isInteger"`
	IsNumber    *bool `json:"isNumber"`
	IsString    *bool `json:"isString"`
	IsObject    *bool `json:"isObject"`
	IsSymbol    *bool `json:"isSymbol"`
	IsBigInt    *bool `json:"isBigInt"`
}

type TypeDescription struct {
	IsValid             *bool                   `json:"isValid"`
	LeastCommonAncestor *string                 `json:"leastCommonAncestor,omitempty"`
	TypeSet             *TypeSet                `json:"typeSet,omitempty"`
	Structures          *[]StructureDescription `json:"structures,omitempty"`
	IsTruncated         *bool                   `json:"isTruncated,omitempty"`
}

type TypeLocation struct {
	TypeInformationDescriptor *int    `json:"typeInformationDescriptor"`
	SourceID                  *string `json:"sourceID"`
	Divot                     *int    `json:"divot"`
}

type BasicBlock struct {
	StartOffset    *int  `json:"startOffset"`
	EndOffset      *int  `json:"endOffset"`
	HasExecuted    *bool `json:"hasExecuted"`
	ExecutionCount *int  `json:"executionCount"`
}
