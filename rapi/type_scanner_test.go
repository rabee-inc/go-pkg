package rapi_test

import (
	"slices"
	"testing"

	"github.com/rabee-inc/go-pkg/rapi"
)

type scanEmbedded struct {
	EmbeddedID string `json:"embedded_id"`
}

type scanTarget struct {
	scanEmbedded
	Name     string         `json:"name" validate:"required"`
	Age      int            `json:"age,omitempty"`
	Rate     float64        `json:"rate"`
	Enabled  bool           `json:"enabled"`
	Tags     []string       `json:"tags"`
	Attrs    map[string]int `json:"attrs"`
	Secret   string         `json:"-"`
	NoTag    string
	Nullable *string `json:"nullable"`
}

func Test_Rapi_TypeScanner_Scan(t *testing.T) {
	type args struct {
		value any
	}
	type want struct {
		kind       string
		goTypeName string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: string",
			args: args{value: "a"},
			want: want{kind: rapi.TypeKindString, goTypeName: "string"},
		},
		{
			name: "正常系: int",
			args: args{value: 1},
			want: want{kind: rapi.TypeKindInt, goTypeName: "int"},
		},
		{
			name: "正常系: int64もintとして扱う",
			args: args{value: int64(1)},
			want: want{kind: rapi.TypeKindInt, goTypeName: "int64"},
		},
		{
			name: "正常系: float64",
			args: args{value: 1.5},
			want: want{kind: rapi.TypeKindFloat, goTypeName: "float64"},
		},
		{
			name: "正常系: bool",
			args: args{value: true},
			want: want{kind: rapi.TypeKindBool, goTypeName: "bool"},
		},
		{
			name: "正常系: slice",
			args: args{value: []string{}},
			want: want{kind: rapi.TypeKindArray, goTypeName: "[]string"},
		},
		{
			name: "正常系: map",
			args: args{value: map[string]int{}},
			want: want{kind: rapi.TypeKindMap, goTypeName: "map[string]int"},
		},
		{
			name: "正常系: pointerは参照先の型として扱う",
			args: args{value: new(string)},
			want: want{kind: rapi.TypeKindString, goTypeName: "string"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := rapi.NewTypeScanner().Scan(tc.args.value)

			if result == nil {
				t.Fatalf("result = nil, want not nil")
			}
			if result.Kind != tc.want.kind {
				t.Errorf("Kind = %v, want %v", result.Kind, tc.want.kind)
			}
			if result.GoTypeName != tc.want.goTypeName {
				t.Errorf("GoTypeName = %v, want %v", result.GoTypeName, tc.want.goTypeName)
			}
		})
	}
}

func Test_Rapi_TypeScanner_Scan_Struct(t *testing.T) {
	type args struct {
		tagNames           []string
		structFieldEnabled bool
	}
	type want struct {
		fieldNames []string
		omitEmpty  map[string]bool
		validate   map[string]string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: jsonタグ有効・タグなしフィールドも含める",
			args: args{tagNames: []string{"json"}, structFieldEnabled: true},
			want: want{
				fieldNames: []string{"NoTag", "age", "attrs", "embedded_id", "enabled", "name", "nullable", "rate", "tags"},
				omitEmpty:  map[string]bool{"age": true, "name": false},
				validate:   map[string]string{"name": "required", "age": ""},
			},
		},
		{
			name: "正常系: jsonタグ有効・タグなしフィールドは除外",
			args: args{tagNames: []string{"json"}, structFieldEnabled: false},
			want: want{
				fieldNames: []string{"age", "attrs", "embedded_id", "enabled", "name", "nullable", "rate", "tags"},
				omitEmpty:  map[string]bool{"age": true, "name": false},
				validate:   map[string]string{"name": "required"},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts := rapi.NewTypeScanner().AddStructTagName(tc.args.tagNames...)
			if tc.args.structFieldEnabled {
				ts = ts.EnableStructField()
			} else {
				ts = ts.DisableStructField()
			}

			result := ts.Scan(scanTarget{})
			if result == nil {
				t.Fatalf("result = nil, want not nil")
			}
			if result.Kind != rapi.TypeKindStruct {
				t.Errorf("Kind = %v, want %v", result.Kind, rapi.TypeKindStruct)
			}

			// Export で embedded がインライン展開される
			exported := ts.Export()
			scanned, ok := exported[result.Name]
			if !ok {
				t.Fatalf("exported[%v] = nil, want not nil", result.Name)
			}

			fieldNames := make([]string, 0, len(scanned.Fields))
			for k := range scanned.Fields {
				fieldNames = append(fieldNames, k)
			}
			slices.Sort(fieldNames)
			if !slices.Equal(fieldNames, tc.want.fieldNames) {
				t.Errorf("fieldNames = %v, want %v", fieldNames, tc.want.fieldNames)
			}

			for k, v := range tc.want.omitEmpty {
				field, ok := scanned.Fields[k]
				if !ok {
					t.Errorf("Fields[%v] = nil, want not nil", k)
					continue
				}
				if field.OmitEmpty != v {
					t.Errorf("Fields[%v].OmitEmpty = %v, want %v", k, field.OmitEmpty, v)
				}
			}
			for k, v := range tc.want.validate {
				field, ok := scanned.Fields[k]
				if !ok {
					t.Errorf("Fields[%v] = nil, want not nil", k)
					continue
				}
				if field.Validate != v {
					t.Errorf("Fields[%v].Validate = %v, want %v", k, field.Validate, v)
				}
			}
		})
	}
}

func Test_Rapi_TypeScanner_Scan_StructFieldKind(t *testing.T) {
	type args struct {
		fieldName string
	}
	type want struct {
		kind         string
		elemKind     string
		keyKind      string
		hasElemType  bool
		hasKeyType   bool
		goTypeNameEq string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: slice は ElemType を持つ",
			args: args{fieldName: "tags"},
			want: want{
				kind:         rapi.TypeKindArray,
				elemKind:     rapi.TypeKindString,
				hasElemType:  true,
				hasKeyType:   false,
				goTypeNameEq: "[]string",
			},
		},
		{
			name: "正常系: map は KeyType と ElemType を持つ",
			args: args{fieldName: "attrs"},
			want: want{
				kind:         rapi.TypeKindMap,
				elemKind:     rapi.TypeKindInt,
				keyKind:      rapi.TypeKindString,
				hasElemType:  true,
				hasKeyType:   true,
				goTypeNameEq: "map[string]int",
			},
		},
		{
			name: "正常系: pointer は参照先の型になる",
			args: args{fieldName: "nullable"},
			want: want{
				kind:         rapi.TypeKindString,
				hasElemType:  false,
				hasKeyType:   false,
				goTypeNameEq: "string",
			},
		},
	}

	ts := rapi.NewTypeScanner().AddStructTagName("json")
	scanned := ts.Scan(scanTarget{})

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			field, ok := scanned.Fields[tc.args.fieldName]
			if !ok {
				t.Fatalf("Fields[%v] = nil, want not nil", tc.args.fieldName)
			}

			if field.Kind != tc.want.kind {
				t.Errorf("Kind = %v, want %v", field.Kind, tc.want.kind)
			}
			if field.GoTypeName != tc.want.goTypeNameEq {
				t.Errorf("GoTypeName = %v, want %v", field.GoTypeName, tc.want.goTypeNameEq)
			}
			if (field.ElemType != nil) != tc.want.hasElemType {
				t.Errorf("ElemType != nil = %v, want %v", field.ElemType != nil, tc.want.hasElemType)
			}
			if (field.KeyType != nil) != tc.want.hasKeyType {
				t.Errorf("KeyType != nil = %v, want %v", field.KeyType != nil, tc.want.hasKeyType)
			}
			if tc.want.hasElemType && field.ElemType.Kind != tc.want.elemKind {
				t.Errorf("ElemType.Kind = %v, want %v", field.ElemType.Kind, tc.want.elemKind)
			}
			if tc.want.hasKeyType && field.KeyType.Kind != tc.want.keyKind {
				t.Errorf("KeyType.Kind = %v, want %v", field.KeyType.Kind, tc.want.keyKind)
			}
		})
	}
}

func Test_Rapi_TypeScanner_Export(t *testing.T) {
	type args struct {
		values []any
	}
	type want struct {
		count int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: scanした構造体がexportされる",
			args: args{values: []any{scanTarget{}}},
			want: want{count: 2}, // scanTarget と scanEmbedded
		},
		{
			name: "正常系: 構造体以外はexportされない",
			args: args{values: []any{"a", 1, []string{}}},
			want: want{count: 0},
		},
		{
			name: "正常系: 同じ型を複数回scanしても1件",
			args: args{values: []any{scanEmbedded{}, scanEmbedded{}}},
			want: want{count: 1},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts := rapi.NewTypeScanner().AddStructTagName("json")
			for _, v := range tc.args.values {
				ts.Scan(v)
			}

			result := ts.Export()

			if len(result) != tc.want.count {
				t.Errorf("len(result) = %v, want %v", len(result), tc.want.count)
			}
		})
	}
}

type unionString string

type unionInt int

func Test_Rapi_TypeScanner_ScanUnion(t *testing.T) {
	type args struct {
		values []any
	}
	type want struct {
		isNil      bool
		name       string
		kind       string
		valueCount int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 文字列のunion",
			args: args{values: []any{unionString("a"), unionString("b")}},
			want: want{isNil: false, name: "unionString", kind: rapi.TypeKindString, valueCount: 2},
		},
		{
			name: "正常系: 整数のunion",
			args: args{values: []any{unionInt(1)}},
			want: want{isNil: false, name: "unionInt", kind: rapi.TypeKindInt, valueCount: 1},
		},
		{
			name: "正常系: 組み込み型はその型名になる",
			args: args{values: []any{"a"}},
			want: want{isNil: false, name: "string", kind: rapi.TypeKindString, valueCount: 1},
		},
		{
			name: "異常系: 空スライスはnil",
			args: args{values: []any{}},
			want: want{isNil: true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := rapi.NewTypeScanner().ScanUnion(tc.args.values)

			if (result == nil) != tc.want.isNil {
				t.Fatalf("result == nil = %v, want %v", result == nil, tc.want.isNil)
			}
			if result == nil {
				return
			}
			if result.Name != tc.want.name {
				t.Errorf("Name = %v, want %v", result.Name, tc.want.name)
			}
			if result.Kind != tc.want.kind {
				t.Errorf("Kind = %v, want %v", result.Kind, tc.want.kind)
			}
			if len(result.Values) != tc.want.valueCount {
				t.Errorf("len(Values) = %v, want %v", len(result.Values), tc.want.valueCount)
			}
		})
	}
}

func Test_Rapi_TypeScanner_ExportUnion(t *testing.T) {
	type args struct {
		values [][]any
	}
	type want struct {
		count int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: scanしたunionがexportされる",
			args: args{values: [][]any{
				{unionString("a")},
				{unionInt(1)},
			}},
			want: want{count: 2},
		},
		{
			name: "正常系: 同じ型を複数回scanしても1件",
			args: args{values: [][]any{
				{unionString("a")},
				{unionString("b")},
			}},
			want: want{count: 1},
		},
		{
			name: "正常系: scanしていない場合は0件",
			args: args{values: [][]any{}},
			want: want{count: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ts := rapi.NewTypeScanner()
			for _, v := range tc.args.values {
				ts.ScanUnion(v)
			}

			result := ts.ExportUnion()

			if len(result) != tc.want.count {
				t.Errorf("len(result) = %v, want %v", len(result), tc.want.count)
			}
		})
	}
}
