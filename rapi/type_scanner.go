package rapi

// 型情報を読み取り、整理して、 JSON 化可能な形で出力する interface
type TypeScanner interface {
	// any の型情報を出力する
	Scan(value any) *TypeStructure
	// 今まで scan したすべての型情報を slice で出力する
	Export() map[string]*TypeStructure
	ScanUnion(values []any) *UnionStructure
	ExportUnion() map[string]*UnionStructure
	EnableStructField() TypeScanner
	DisableStructField() TypeScanner
	AddStructTagName(tagName ...string) TypeScanner
}

// 一つの型情報
type TypeStructure struct {
	Name       string `json:"name"`
	GoTypeName string `json:"go_type_name"`
	Kind       string `json:"kind"`
	// map の key の型情報。map じゃない場合は nil
	KeyType *TypeStructure `json:"key_type,omitempty"`
	// map, slice, array の要素の型情報。それ以外は nil
	ElemType *TypeStructure `json:"elem_type,omitempty"`
	// struct の field の型情報。それ以外は nil
	Fields    map[string]*TypeStructure `json:"fields,omitempty"`
	OmitEmpty bool                      `json:"omit_empty,omitempty"`
	Validate  string                    `json:"validate,omitempty"`
}

type UnionStructure struct {
	Name       string `json:"name"`
	GoTypeName string `json:"go_type_name"`
	Kind       string `json:"kind"`
	Values     []any  `json:"values"`
}

const (
	TypeKindString = "string"
	TypeKindInt    = "int"
	TypeKindFloat  = "float"
	TypeKindBool   = "bool"
	TypeKindArray  = "array"
	TypeKindMap    = "map"
	TypeKindStruct = "struct"
	TypeKindAny    = "any"
)
