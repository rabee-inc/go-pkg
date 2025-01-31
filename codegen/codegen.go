package codegen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rabee-inc/go-pkg/sliceutil"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v3"
)

const (
	typeString      = "string"
	typeStringSlice = "[]string"
	typeInt         = "int"
	typeIntSlice    = "[]int"
	typeFloat       = "float"
	typeFloatSlice  = "[]float"
	typeInt64       = "int64"
	typeInt64Slice  = "[]int64"
)

var reFloat = regexp.MustCompile(`float$`)

func actualType(t string) string {
	// 末尾が float のものは float64 に変換
	if reFloat.MatchString(t) {
		return reFloat.ReplaceAllString(t, "float64")
	}
	return t
}

var vl = validator.New()

// ExportByYaml ... yaml ファイルからコードを生成し、ファイルに出力する
func ExportByYaml(path string) string {
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
	return absOutput
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
		if td.BaseType == typeString {
			outputCode += formatConstantMethodString(pascalName)
		}
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
				var param string
				if extend.IsSlice {
					param = formatConstantMetaDataSliceParam(extend.Name, extend.Type, extend.SliceValue, extend.HasDoubleQuote)
				} else {
					param = formatConstantMetaDataParam(extend.Name, extend.Value, extend.HasDoubleQuote)
				}
				params = append(params, param)
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
