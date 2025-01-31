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

		if n, ok := def.OtherProps.(string); ok {
			it.Name = n
		} else if m, ok := def.OtherProps.(map[string]any); ok {
			// id
			if mid, ok := m["id"]; ok {
				if _, ok := mid.(string); !ok {
					panic(fmt.Sprintf("%s (%s) invalid def: id must be string. (Even if the type is numeric, it must be specified as a string.)", typeName, variableName))
				}
				it.VariableValue = mid.(string)
			}

			// name
			if mName, ok := m["name"]; ok {
				if _, ok := mName.(string); !ok {
					panic(fmt.Sprintf("%s (%s) invalid def: name must be string.", typeName, variableName))
				}
				it.Name = mName.(string)
			} else {
				panic(fmt.Sprintf("%s (%s) invalid def: name is required.", typeName, variableName))
			}

			// extends
			for _, ex := range td.Extends {
				if value, ok := m[ex.Name]; ok {
					hasDQ := ex.Type == typeString || ex.Type == typeStringSlice
					isSlice := strings.HasPrefix(ex.Type, "[]")

					// slice の場合
					if vSlice, ok := value.([]any); ok {
						sliceValue := []string{}
						for _, v := range vSlice {
							if _, ok := v.(string); !ok {
								panic(fmt.Sprintf("%s (%s) invalid def: %s must be string. (Even if the type is numeric, it must be specified as a string.)", typeName, variableName, ex.Name))
							}
							sliceValue = append(sliceValue, v.(string))
						}

						it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
							Name:           ex.Name,
							Type:           ex.Type,
							SliceValue:     sliceValue,
							HasDoubleQuote: hasDQ,
							IsSlice:        isSlice,
						})
					} else {
						// slice 以外
						if _, ok := value.(string); !ok {
							panic(fmt.Sprintf("%s (%s) invalid def: %s must be string. (Even if the type is numeric, it must be specified as a string.)", typeName, variableName, ex.Name))
						}

						it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
							Name:           ex.Name,
							Type:           ex.Type,
							Value:          value.(string),
							HasDoubleQuote: hasDQ,
							IsSlice:        isSlice,
						})
					}
				} else {
					panic(fmt.Sprintf("%s (%s) invalid def: %s (by extends) is required.", typeName, variableName, ex.Name))
				}
			}
		} else {
			panic(fmt.Sprintf("%s (%s) invalid def:\n=== def format ===\n\nid_value_and_variable_name: name_text\n\nor\n\nid_value_and_variable_name:\n  name: name_text\n  prop1: value1\n\nor\n\nvariable_name:\n  id: id_value\n  name: name_text\n", typeName, variableName))
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
