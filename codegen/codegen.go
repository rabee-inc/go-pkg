package codegen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/rabee-inc/go-pkg/sliceutil"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v3"
)

const (
	typeString = "string"
	typeInt    = "int"
	typeFloat  = "float"
	typeInt64  = "int64"
)

func actualType(t string) string {
	if t == typeFloat {
		return "float64"
	}
	return t
}

var vl = validator.New()

// ExportByYaml ... yaml ファイルからコードを生成し、ファイルに出力する
func ExportByYaml(path string) {
	// yaml 読み込み
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	formattedCode, val := GenerateByYamlFile(filepath.Base(path), file)

	absOutput, err := filepath.Abs(val.Settings.Output)
	if err != nil {
		panic(err)
	}

	// ディレクトリがない場合は、作る
	dir := filepath.Dir(absOutput)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(absOutput, formattedCode, 0666)
	if err != nil {
		panic(err)
	}

	fmt.Println("generated: " + absOutput)
}

// GenerateByYamlFile ... yaml ファイルからコードを生成する
func GenerateByYamlFile(name string, file []byte) ([]byte, *yamlInput) {
	val := &yamlInput{}
	err := yaml.Unmarshal(file, &val)
	if err != nil {
		panic(err)
	}

	// validate
	if err := vl.Struct(val); err != nil {
		panic(err)
	}

	typeDefs := []*typeDef{}

	for _, v := range val.Types {
		typeDefs = append(typeDefs, newTypeDef(v))
	}

	outputCode := formatHeader(name) + "\n\n"
	outputCode += formatPackage(val.Settings.Package) + "\n\n"
	outputCode += formatCheckSum(GenerateCheckSum(file)) + "\n\n"
	outputCode += defaultMetaDataCode + "\n\n"

	// constants struct
	constantsStructParams := []string{}

	// generate any slice for GetConstIDs
	generateAnySliceCodes := []string{}
	anySliceVars := []string{}

	// init map generate codes
	generateMapCodes := []string{}

	// init constants params
	constantsInitParams := []string{}

	for _, td := range typeDefs {
		pascalName := toPascalCase(td.Name)
		// comment
		outputCode += formatConstantComment(pascalName, td.Comment)
		// type
		outputCode += formatConstantType(pascalName, td.BaseType)

		// method
		outputCode += formatConstantMethodString(pascalName)
		outputCode += formatConstantMethodMeta(pascalName)
		outputCode += formatConstantMethodName(pascalName)

		// const
		constValues := []string{}
		for _, def := range td.Defs {
			constValues = append(constValues, formatConstantValue(pascalName, toPascalCase(def.VariableName), td.BaseType == typeString, def.VariableValue))
		}
		outputCode += formatConstantValues(strings.Join(constValues, "\n"))

		// meta data type
		if td.HasExtends {
			params := []string{}
			params = append(params, formatConstantMetaDataTypeID(pascalName))
			params = append(params, formatConstantMetaDataTypeParam("name", "string"))
			for _, def := range td.Extends {
				params = append(params, formatConstantMetaDataTypeParam(def.Name, def.Type))
			}
			outputCode += formatConstantMetaDataType(pascalName, strings.Join(params, "\n"))
		} else {
			outputCode += formatConstantMetaDataByGenerics(pascalName)
		}

		// meta data list
		metaDataListElements := []string{}
		for _, def := range td.Defs {
			params := []string{}
			params = append(params, formatConstantMetaDataParam("ID", pascalName+toPascalCase(def.VariableName), false))
			params = append(params, formatConstantMetaDataParam("Name", def.Name, true))
			for _, extend := range def.ExtendValues {
				params = append(params, formatConstantMetaDataParam(extend.Name, extend.Value, extend.HasDoubleQuote))
			}
			metaDataListElements = append(metaDataListElements, formatConstantMetaDataListElement(strings.Join(params, "\n")))
		}
		outputCode += formatConstantMetaDataList(pascalName, strings.Join(metaDataListElements, "\n"))

		// meta data map (var)
		outputCode += formatConstantMetaDataMap(toPascalCase(td.Name))

		if !td.OnlyBackend {
			// constant struct params
			constantsStructParams = append(constantsStructParams, formatConstantsStructParam(toPluralForm(td.Name), "[]*"+formatConstantMetaDataTypeName(pascalName)))
			constantsStructParams = append(constantsStructParams, formatConstantsStructParam(td.Name, formatConstantMetaDataMapTypeName(pascalName)))

			// init map generate codes
			generateMapCodes = append(generateMapCodes, formatGenerateMapCode(pascalName))

			// generate any slice for GetConstIDs
			anySliceVarName := strings.ToLower(pascalName[0:1]) + pascalName[1:]
			generateAnySliceCodes = append(generateAnySliceCodes, formatGenerateAnySliceCode(anySliceVarName, pascalName))
			anySliceVars = append(anySliceVars, anySliceVarName+",")

			// init constants params
			constantsInitParams = append(constantsInitParams, formatConstantsInitParam(toPluralForm(pascalName), toPluralForm(pascalName)))
			constantsInitParams = append(constantsInitParams, formatConstantsInitParam(pascalName, formatConstantMetaDataMapVariableName(pascalName)))
		}
	}

	outputCode += formatConstantsStruct(strings.Join(constantsStructParams, "\n"))
	outputCode += formatConstantsMethodGetConstIDs(strings.Join(generateAnySliceCodes, "\n"), strings.Join(anySliceVars, "\n"))
	outputCode += formatInitCode(strings.Join(generateMapCodes, "\n"), strings.Join(constantsInitParams, "\n"))

	// コードのフォーマット
	formattedCode, err := format.Source([]byte(outputCode))
	if err != nil {
		panic(err)
	}

	return formattedCode, val
}

// GenerateCheckSum ... チェックサムを生成する
func GenerateCheckSum(text []byte) string {
	r := sha256.Sum256(text)
	checkSum := hex.EncodeToString(r[:])
	return checkSum
}

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
			for _, v := range td.Extends {
				if value, ok := m[v.Name]; ok {
					if _, ok := value.(string); !ok {
						panic(fmt.Sprintf("%s (%s) invalid def: %s must be string. (Even if the type is numeric, it must be specified as a string.)", typeName, variableName, v.Name))
					}
					it.ExtendValues = append(it.ExtendValues, &metaDataValueDef{
						Name:           v.Name,
						Value:          value.(string),
						HasDoubleQuote: v.Type == typeString,
					})
				} else {
					panic(fmt.Sprintf("%s (%s) invalid def: %s (by extends) is required.", typeName, variableName, v.Name))
				}
			}
		} else {
			panic(fmt.Sprintf("%s (%s) invalid def:\n=== def format ===\n\nid_value_and_variable_name: name_text\n\nor\n\nid_value_and_variable_name:\n  name: name_text\n  prop1: value1\n\nor\n\nvariable_name:\n  id: id_value\n  name: name_text\n", typeName, variableName))
		}
	}

	return td
}

type yamlInput struct {
	Settings *settingsInput `yaml:"settings" validate:"required"`
	Types    typeInputList  `yaml:"types" validate:"required"`
}

type typeInputList []*typeInput

func (p *typeInputList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`types` must contain YAML mapping, has %v", value.Kind)
	}

	*p = make([]*typeInput, len(value.Content)/2)

	for i := 0; i < len(value.Content); i += 2 {
		ti := &typeInput{}
		var typeName string

		// key
		if err := value.Content[i].Decode(&typeName); err != nil {
			return err
		}
		// value
		if err := value.Content[i+1].Decode(&ti); err != nil {
			return err
		}

		ti.Name = typeName
		(*p)[i/2] = ti
	}
	return nil
}

type settingsInput struct {
	Package string `yaml:"package" validate:"required"`
	Output  string `yaml:"output" validate:"required"`
}

type defInput struct {
	Name       string
	OtherProps any
}

type defInputList []*defInput

func (p *defInputList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`defs` must contain YAML mapping, has %v", value.Kind)
	}

	*p = make([]*defInput, len(value.Content)/2)

	for i := 0; i < len(value.Content); i += 2 {
		di := &defInput{}
		var keyName string

		// key
		if err := value.Content[i].Decode(&keyName); err != nil {
			return err
		}
		// value
		if err := value.Content[i+1].Decode(&di.OtherProps); err != nil {
			return err
		}

		di.Name = keyName
		(*p)[i/2] = di
	}
	return nil
}

type extendsInputList []*extendDef

func (p *extendsInputList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`types` must contain YAML mapping, has %v", value.Kind)
	}

	*p = make([]*extendDef, len(value.Content)/2)

	for i := 0; i < len(value.Content); i += 2 {
		var exName string
		var typeName string

		// key
		if err := value.Content[i].Decode(&exName); err != nil {
			return err
		}
		// value
		if err := value.Content[i+1].Decode(&typeName); err != nil {
			return err
		}

		(*p)[i/2] = &extendDef{
			Name: exName,
			Type: typeName,
		}
	}
	return nil
}

type typeInput struct {
	Name        string
	Comment     string           `yaml:"comment" validate:"required"`
	Type        string           `yaml:"type"`
	OnlyBackend bool             `yaml:"only_backend"`
	Extends     extendsInputList `yaml:"extends"`
	Defs        defInputList     `yaml:"defs" validate:"required"`
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
	Value          string
	HasDoubleQuote bool
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

// hello world → HelloWorld
// hello_world → HelloWorld
// hello-world → HelloWorld
// hello → Hello
func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	ss := strings.Split(s, " ")
	res := ""
	for _, str := range ss {
		if len(str) > 0 {
			res += strings.ToUpper(str[0:1]) + str[1:]
		}
	}
	return res
}

func toPluralForm(word string) string {
	if strings.HasSuffix(word, "s") {
		return word + "es"
	} else if strings.HasSuffix(word, "y") {
		// (a,i,u,e,o)y で終わる場合、sだけつける
		if sliceutil.Contains([]string{"a", "i", "u", "e", "o"}, string(word[len(word)-2])) {
			return word + "s"
		}
		return word[:len(word)-1] + "ies"
	} else {
		return word + "s"
	}
}

const header = "// Code generated by %s DO NOT EDIT."

func formatHeader(name string) string {
	return fmt.Sprintf(header, name)
}

const packageCode = "package %s"

func formatPackage(name string) string {
	return fmt.Sprintf(packageCode, name)
}

const checkSumCode = `const CheckSum = "%s"`

func formatCheckSum(checkSum string) string {
	return fmt.Sprintf(checkSumCode, checkSum)
}

const defaultMetaDataCode = `
type ConstantMetaData[T comparable] struct {
	ID   T      ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

`

const constantCommentCode = "// %s ... %s"

func formatConstantComment(name, comment string) string {
	return fmt.Sprintf(constantCommentCode, name, comment)
}

const constantTypeCode = `
type %s %s

`

func formatConstantType(name, base string) string {
	return fmt.Sprintf(constantTypeCode, name, base)
}

const constantValuesCode = `
const (
%s
)

`

const constantMethodStringCode = `
func (c %s) String() string {
	return string(c)
}

`

func formatConstantMethodString(name string) string {
	return fmt.Sprintf(constantMethodStringCode, name)
}

const constantMethodMetaCode = `
func (c %s) Meta() (*%s, bool) {
	m, ok := %s[c]
	return m, ok
}

`

func formatConstantMethodMeta(name string) string {
	tName := formatConstantMetaDataTypeName(name)
	mapName := formatConstantMetaDataMapVariableName(name)
	return fmt.Sprintf(constantMethodMetaCode, name, tName, mapName)
}

const constantMethodNameCode = `
func (c %s) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

`

func formatConstantMethodName(name string) string {
	return fmt.Sprintf(constantMethodNameCode, name)
}

func formatConstantValues(values string) string {
	return fmt.Sprintf(constantValuesCode, values)
}

const constantValueCode = `%s%s %s = %s`

func formatConstantValue(tName, vName string, hasDoubleQuote bool, value string) string {
	outputValue := value
	if hasDoubleQuote {
		outputValue = fmt.Sprintf(`"%s"`, value)
	}
	return fmt.Sprintf(constantValueCode, tName, vName, tName, outputValue)
}

const constantMetaDataTypeNameCode = `%sMetaData`

func formatConstantMetaDataTypeName(tName string) string {
	return fmt.Sprintf(constantMetaDataTypeNameCode, tName)
}

const constantMetaDataTypeByGenericsCode = `
type %s ConstantMetaData[%s]

`

func formatConstantMetaDataByGenerics(tName string) string {
	return fmt.Sprintf(constantMetaDataTypeByGenericsCode, formatConstantMetaDataTypeName(tName), tName)
}

const constantMetaDataTypeCode = `
type %s struct {
%s
}

`

func formatConstantMetaDataType(tName, params string) string {
	return fmt.Sprintf(constantMetaDataTypeCode, formatConstantMetaDataTypeName(tName), params)
}

const constantMetaDataTypeCodeParam = `%s   %s ` + "`json:\"%s\"`"

func formatConstantMetaDataTypeParam(name, tName string) string {
	return fmt.Sprintf(constantMetaDataTypeCodeParam, toPascalCase(name), tName, name)
}

func formatConstantMetaDataTypeID(tName string) string {
	return fmt.Sprintf(constantMetaDataTypeCodeParam, "ID", tName, "id")
}

const constantMetaDataMapVariableNameCode = `%sMap`

func formatConstantMetaDataMapVariableName(tName string) string {
	return fmt.Sprintf(constantMetaDataMapVariableNameCode, tName)
}

const constantMetaDataMapTypeCode = `map[%s]*%s`

func formatConstantMetaDataMapTypeName(tName string) string {
	return fmt.Sprintf(constantMetaDataMapTypeCode, tName, formatConstantMetaDataTypeName(tName))
}

const constantMetaDataMapCode = `
var %s %s

`

func formatConstantMetaDataMap(tName string) string {
	return fmt.Sprintf(constantMetaDataMapCode, formatConstantMetaDataMapVariableName(tName), formatConstantMetaDataMapTypeName(tName))
}

const constantMetaDataListCode = `
var %s = []*%s{
%s
}

`

func formatConstantMetaDataList(tName, elements string) string {
	return fmt.Sprintf(constantMetaDataListCode, toPluralForm(tName), formatConstantMetaDataTypeName(tName), elements)
}

const constantMetaDataListElementCode = `{
%s
},`

func formatConstantMetaDataListElement(params string) string {
	return fmt.Sprintf(constantMetaDataListElementCode, params)
}

const constantMetaDataParamCode = `%s: %s,`

func formatConstantMetaDataParam(name, value string, hasDoubleQuote bool) string {
	if hasDoubleQuote {
		value = fmt.Sprintf(`"%s"`, value)
	}
	return fmt.Sprintf(constantMetaDataParamCode, toPascalCase(name), value)
}

const constantsTypeNameCode = `Constants`
const constantsVariableNameCode = `ConstantsData`

const constantsStructCode = `
type ` + constantsTypeNameCode + ` struct {
%s
}

var ` + constantsVariableNameCode + ` *` + constantsTypeNameCode + `

`

func formatConstantsStruct(params string) string {
	return fmt.Sprintf(constantsStructCode, params)
}

const constantsStructParamCode = `%s %s ` + "`json:\"%s\"`"

func formatConstantsStructParam(name, tName string) string {
	return fmt.Sprintf(constantsStructParamCode, toPascalCase(name), tName, name)
}

const constantsMethodGetConstIDsCode = `
// deprecated use ConstIDs
func (c *` + constantsTypeNameCode + `) GetConstIDs() [][]any {
	%s
	return [][]any{
		%s
	}
}

func (c *` + constantsTypeNameCode + `) ConstIDs() [][]any {
	return c.GetConstIDs()
}

`

func formatConstantsMethodGetConstIDs(generateAnySliceCodes, constantsParams string) string {
	return fmt.Sprintf(constantsMethodGetConstIDsCode, generateAnySliceCodes, constantsParams)
}

const generateAnySliceCode = `
%s := []any{}
for _, v := range c.%s {
	%s = append(%s, v.ID)
}
`

func formatGenerateAnySliceCode(name, tName string) string {
	return fmt.Sprintf(generateAnySliceCode, name, toPluralForm(tName), name, name)
}

const initCode = `
func init() {
	%s
	` + constantsVariableNameCode + ` = &` + constantsTypeNameCode + `{
		%s
	}
}
`

func formatInitCode(generateMapCodes, constantsParams string) string {
	return fmt.Sprintf(initCode, generateMapCodes, constantsParams)
}

const generateMapCode = `
	%s = %s{}
	for _, v := range %s {
		%s[v.ID] = v
	}
`

func formatGenerateMapCode(tName string) string {
	return fmt.Sprintf(generateMapCode,
		formatConstantMetaDataMapVariableName(tName),
		formatConstantMetaDataMapTypeName(tName),
		toPluralForm(tName),
		formatConstantMetaDataMapVariableName(tName),
	)
}

const constantsInitParamCode = `%s: %s,`

func formatConstantsInitParam(name, value string) string {
	return fmt.Sprintf(constantsInitParamCode, name, value)
}
