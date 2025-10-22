package rapi

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

type typeScanner struct {
	types              map[string]*TypeStructure
	unions             map[string]*UnionStructure
	structFieldEnabled bool
	structTagNames     []string
	unnamedCount       int
}

func NewTypeScanner() TypeScanner {
	return &typeScanner{
		types:              map[string]*TypeStructure{},
		unions:             map[string]*UnionStructure{},
		structFieldEnabled: true,
		structTagNames:     []string{},
	}
}

func (t *typeScanner) EnableStructField() TypeScanner {
	t.structFieldEnabled = true
	return t
}

func (t *typeScanner) DisableStructField() TypeScanner {
	t.structFieldEnabled = false
	return t
}

func (t *typeScanner) AddStructTagName(tagName ...string) TypeScanner {
	t.structTagNames = append(t.structTagNames, tagName...)
	return t
}

func (t *typeScanner) Scan(value any) *TypeStructure {
	return t.scan(reflect.TypeOf(value), false)
}

func (ts *TypeStructure) getFieldsRemovedStruct() *TypeStructure {
	copied := &TypeStructure{}
	copied.Name = ts.Name
	copied.GoTypeName = ts.GoTypeName
	copied.Kind = ts.Kind
	return copied
}

func (t *typeScanner) scan(rt reflect.Type, ignoreField bool) *TypeStructure {
	if rt == nil {
		return nil
	}

	// pointer の場合は pointer じゃなくなるまで Elem
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var typeKind string

	switch rt.Kind() {
	case reflect.String:
		typeKind = TypeKindString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeKind = TypeKindInt
	case reflect.Float32, reflect.Float64:
		typeKind = TypeKindFloat
	case reflect.Bool:
		typeKind = TypeKindBool
	case reflect.Map:
		typeKind = TypeKindMap
	case reflect.Slice, reflect.Array, reflect.Chan:
		typeKind = TypeKindArray
	case reflect.Struct:
		typeKind = TypeKindStruct
	case reflect.Interface:
		typeKind = TypeKindAny
	default:
		panic(fmt.Sprintf("Invalid type: %s. (%s)", rt.Kind().String(), rt.String()))
	}

	ts := &TypeStructure{
		Kind: typeKind,
	}

	switch rt.Kind() {
	// map
	case reflect.Map:
		ts.Name = rt.Name()
		ts.GoTypeName = rt.String()
		ts.KeyType = t.scan(rt.Key(), true)
		ts.ElemType = t.scan(rt.Elem(), true)

	// array
	case reflect.Slice, reflect.Array, reflect.Chan:
		ts.Name = rt.Name()
		ts.GoTypeName = rt.String()
		ts.ElemType = t.scan(rt.Elem(), true)

	// struct
	case reflect.Struct:
		name := rt.Name()

		if name == "" {
			name = "__unnamed__." + string(rune(t.unnamedCount))
			t.unnamedCount++
		} else {
			name = rt.String()
		}

		if v, ok := t.types[name]; ok {
			if ignoreField {
				return v.getFieldsRemovedStruct()
			}
			return v
		}

		ts.Name = name
		ts.GoTypeName = rt.String()
		ts.Fields = map[string]*TypeStructure{}

		t.types[name] = ts
		for i := 0; i < rt.NumField(); i++ {
			keyName := ""
			field := rt.Field(i)

			hasSkip := false
			omitEmpty := false
			for _, tagName := range t.structTagNames {
				tagValue := field.Tag.Get(tagName)
				values := strings.Split(tagValue, ",")
				tagValue = values[0]

				if tagValue == "-" {
					hasSkip = true
					break
				}
				keyName = tagValue
				if keyName != "" {
					omitEmpty = slices.Contains(values[1:], "omitempty")
					break
				}
			}

			if hasSkip {
				continue
			}

			if keyName == "" {
				// embedded field の場合は、自身のフィールドとして処理する
				if field.Anonymous {
					fieldTs := t.scan(field.Type, false)
					if fieldTs != nil {
						if ts.InlineEmbeddedFields == nil {
							ts.InlineEmbeddedFields = map[string]*TypeStructure{}
						}
						// scan 完了後にインライン展開させるために登録しておく
						ts.InlineEmbeddedFields[fieldTs.Name] = fieldTs
					}
					continue
				}

				// struct tag による命名がなく、 structFieldEnabled が false の場合は、そのフィールドは存在しないものとして扱う
				if !t.structFieldEnabled {
					continue
				}
				// tag による命名がない場合は、フィールド名をそのまま使用する
				keyName = field.Name
			}

			fieldTs := t.scan(field.Type, true)
			if fieldTs != nil {
				fieldTs.OmitEmpty = omitEmpty
				fieldTs.Validate = field.Tag.Get("validate")
				ts.Fields[keyName] = fieldTs
			}
		}

		// json 化するときに再帰的に参照され続けないように fields を削除
		if ignoreField {
			return ts.getFieldsRemovedStruct()
		}

	// primitive or other
	default:
		ts.Name = rt.Name()
		ts.GoTypeName = rt.String()
	}

	return ts
}

func (t *typeScanner) formatScannedTypeStructure(ts *TypeStructure) {
	for _, embeddedTs := range ts.InlineEmbeddedFields {
		for k, v := range embeddedTs.Fields {
			ts.Fields[k] = v
		}
	}

	for _, field := range ts.Fields {
		t.formatScannedTypeStructure(field)
	}
}

func (t *typeScanner) Export() map[string]*TypeStructure {
	// types をコピーして返す
	types := map[string]*TypeStructure{}
	for k, v := range t.types {
		t.formatScannedTypeStructure(v)
		types[k] = v
	}
	return types
}

func (t *typeScanner) ScanUnion(values []any) *UnionStructure {
	if len(values) == 0 {
		return nil
	}
	rt := reflect.TypeOf(values[0])
	typeName := rt.Name()
	if typeName == "" {
		return nil
	}

	var kind string
	switch rt.Kind() {
	case reflect.String:
		kind = TypeKindString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		kind = TypeKindInt
	case reflect.Float32, reflect.Float64:
		kind = TypeKindFloat
	default:
		panic(fmt.Sprintf("Invalid type: %s. (%s)", rt.Kind().String(), rt.String()))
	}

	us := &UnionStructure{
		Name:       typeName,
		GoTypeName: rt.String(),
		Kind:       kind,
		Values:     values,
	}

	t.unions[us.GoTypeName] = us
	return us
}

func (t *typeScanner) ExportUnion() map[string]*UnionStructure {
	// unions をコピーして返す
	unions := map[string]*UnionStructure{}
	for k, v := range t.unions {
		unions[k] = v
	}
	return unions
}
