package codegen

import (
	"fmt"
	"strings"
)

func newTypeDef(ts *typeInput) *typeDef {
	typeName := ts.Name

	td := &typeDef{
		Name:        typeName,
		Extends:     []*extendDef{},
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
		for _, v := range ts.Extends {
			td.Extends = append(td.Extends, &extendDef{
				Name: v.Name,
				Type: actualType(v.Type),
			})
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

					it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
						Name:           ex.Name,
						Type:           ex.Type,
						SliceValue:     value.Values,
						HasDoubleQuote: hasDQ,
						IsSlice:        true,
					})
				} else {
					// slice 以外
					if value.IsSlice {
						panic(fmt.Sprintf("%s (%s) invalid def: %s must not be slice.", typeName, variableName, ex.Name))
					}

					it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
						Name:           ex.Name,
						Type:           ex.Type,
						Value:          value.Value,
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

type extendDef struct {
	Name string
	Type string
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
	Extends     []*extendDef
	Defs        []*typeDefsItem
}
