package codegen

import (
	"fmt"
	"strings"
)

// yamlの入力値から実際のgo型に変換
func actualType(t string) string {
	switch t {
	case typeFloat, typeFloatSlice:
		// float64 に変換
		return strings.Replace(t, "float", "float64", 1)
	default:
		return t
	}
}

// slice の場合でも正しく pascal case に変換する
func toPascalCaseType(t string) string {
	if strings.HasPrefix(t, "[]") {
		return "[]" + toPascalCase(strings.TrimPrefix(t, "[]"))
	}
	return toPascalCase(t)
}

// extendDef(extends 内の 1プロパティ) を生成する
func newExtendDef(ex *extendPropInput) *extendPropDef {
	isPrimitive := primitiveTypeSet.Has(ex.Type)
	var t string
	if isPrimitive {
		t = actualType(ex.Type)
	} else {
		// ユーザーが定義した型を参照した場合は pascal case に変換
		t = toPascalCaseType(ex.Type)
	}
	return &extendPropDef{
		Name:        ex.Name,
		IsPrimitive: isPrimitive,
		Type:        t,
	}
}

// extendsDef を生成する
func newExtendsDef(exDef *extendsDefInput) *extendsDef {
	extends := []*extendPropDef{}
	for _, ex := range exDef.Extends {
		extends = append(extends, newExtendDef(ex))
	}
	return &extendsDef{
		Name:    exDef.Name,
		Extends: extends,
	}
}

// typeDef を生成する
func newTypeDef(ts *typeInput) *typeDef {
	typeName := ts.Name

	td := &typeDef{
		Name:        typeName,
		Extends:     []*extendPropDef{},
		Defs:        []*typeDefsItem{},
		BaseType:    typeString,
		OnlyBackend: ts.OnlyBackend,
	}

	// Comment
	if ts.Comment == "" {
		panic(fmt.Sprintf("%s comment is required.", typeName))
	}
	td.Comment = ts.Comment

	// BaseType
	if ts.Type != "" {
		switch ts.Type {
		case typeString, typeInt, typeFloat, typeInt64:
			td.BaseType = actualType(ts.Type)
		default:
			types := strings.Join([]string{typeString, typeInt, typeFloat, typeInt64}, ", ")
			panic(fmt.Sprintf("%s invalid type: %s.\nAvailable types: %s", typeName, ts.Type, types))
		}
	}

	// Extends
	if len(ts.Extends) > 0 {
		td.HasExtends = true
		for _, ex := range ts.Extends {
			td.Extends = append(td.Extends, newExtendDef(ex))
		}
	}

	// Defs
	for _, def := range ts.Defs {
		variableName := def.Name
		it := &typeDefsItem{
			VariableName:  variableName,
			VariableValue: variableName,
			ExtendValues:  []*metaDataValueDef{},
		}
		td.Defs = append(td.Defs, it)

		// id が指定されている場合上書きする
		if id, ok := def.PropMap["id"]; ok {
			if id.IsSlice {
				panic(fmt.Sprintf("%s (%s) invalid def: id must be string.", typeName, variableName))
			}
			it.VariableValue = id.Value
		}

		// name の設定
		if name, ok := def.PropMap["name"]; ok {
			if name.IsSlice {
				panic(fmt.Sprintf("%s (%s) invalid def: name must be string.", typeName, variableName))
			}
			it.Name = name.Value
		} else {
			// name は必須
			panic(fmt.Sprintf("%s (%s) invalid def: name is required.", typeName, variableName))
		}

		// extends
		for _, ex := range td.Extends {
			if value, ok := def.PropMap[ex.Name]; ok {
				hasDQ := ex.Type == typeString || ex.Type == typeStringSlice
				isSlice := strings.HasPrefix(ex.Type, "[]")

				// slice の場合
				if isSlice {
					if !value.IsSlice {
						panic(fmt.Sprintf("%s (%s) invalid def: %s must be slice.", typeName, variableName, ex.Name))
					}

					sliceValue := value.Values
					if !ex.IsPrimitive {
						tmp := make([]string, len(sliceValue))
						// ユーザー定義の型の場合は 型名と値を pascal case に変換したsliceにする
						for i, v := range sliceValue {
							tmp[i] = strings.TrimPrefix(ex.Type, "[]") + toPascalCase(v)
						}
						sliceValue = tmp
					}
					it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
						Name:           ex.Name,
						Type:           ex.Type,
						SliceValue:     sliceValue,
						HasDoubleQuote: hasDQ,
						IsSlice:        true,
					})
				} else {
					// slice 以外
					if value.IsSlice {
						panic(fmt.Sprintf("%s (%s) invalid def: %s must not be slice.", typeName, variableName, ex.Name))
					}
					v := value.Value
					if !ex.IsPrimitive {
						// ユーザー定義の型の場合は 型名と値を pascal case に変換した値にする
						v = ex.Type + toPascalCase(v)
					}
					it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
						Name:           ex.Name,
						Type:           ex.Type,
						Value:          v,
						HasDoubleQuote: hasDQ,
						IsSlice:        false,
					})
				}
			} else {
				panic(fmt.Sprintf("%s (%s) invalid def: %s (by extends) is required.", typeName, variableName, ex.Name))
			}
		}
	}

	return td
}

type extendPropDef struct {
	Name        string
	IsPrimitive bool
	Type        string
}

type extendsDef struct {
	Name    string
	Extends []*extendPropDef
}

type typeDefsItem struct {
	VariableName  string
	VariableValue string
	Name          string
	ExtendValues  []*metaDataValueDef
}

type metaDataValueDef struct {
	Name           string
	Type           string
	Value          string
	SliceValue     []string
	HasDoubleQuote bool
	IsSlice        bool
}

type typeDef struct {
	Name        string
	Comment     string
	BaseType    string
	OnlyBackend bool
	HasExtends  bool
	Extends     []*extendPropDef
	Defs        []*typeDefsItem
}
