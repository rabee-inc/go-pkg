package codegen

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// このファイルでは yaml からの入力についてのみ処理を行う

// yamlInput ... yaml からの入力を構造化したもの
type yamlInput struct {
	Settings  *settingsInput  `yaml:"settings" validate:"required"`
	Templates *templatesInput `yaml:"templates"`
	Types     typeInputList   `yaml:"types" validate:"required"`
}

// settingsInput ... settings を構造化したもの
type settingsInput struct {
	Package string `yaml:"package" validate:"required"`
	Output  string `yaml:"output" validate:"required"`
}

// templatesInput ... templates を構造化したもの
type templatesInput struct {
	ExtendsDefs extendsDefInputList `yaml:"extends_defs"`
}

// typeInput ... types のそれぞれを構造化したもの
type typeInput struct {
	Name        string
	Comment     string       `yaml:"comment" validate:"required"`
	Type        string       `yaml:"type"`
	OnlyBackend bool         `yaml:"only_backend"`
	Extends     extendsInput `yaml:"extends"`
	Defs        defInputList `yaml:"defs" validate:"required"`
}

// typeInputList ... types を構造化したもの
type typeInputList []*typeInput

// types: の部分
func (p *typeInputList) UnmarshalYAML(value *yaml.Node) error {
	// MappingNode のみ許可
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`types` must contain YAML mapping.")
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

// defInput ... defs のそれぞれを構造化したもの
type defInput struct {
	Name    string
	PropMap defPropInputMap
}

// defInputList ... defs を構造化したもの
type defInputList []*defInput

// defs: の部分
func (p *defInputList) UnmarshalYAML(value *yaml.Node) error {
	// MappingNode のみ許可
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf(errorInvalidKind+"\n"+errorInvalidDefs, value.Kind)
	}

	*p = make([]*defInput, len(value.Content)/2)

	for i := 0; i < len(value.Content); i += 2 {
		di := &defInput{}
		var keyName string

		// key
		if err := value.Content[i].Decode(&keyName); err != nil {
			return err
		}
		// props
		if err := value.Content[i+1].Decode(&di.PropMap); err != nil {
			return err
		}

		di.Name = keyName
		(*p)[i/2] = di
	}
	return nil
}

// defPropInput ... def の各プロパティ
type defPropInput struct {
	Name  string
	Value defPropValueInput
}

// defPropInputMap ... defPropInput の map
type defPropInputMap map[string]defPropValueInput

// 以下の部分
//
//	defs:
//		def_name: xxx
//
//	defs:
//		def_name:
//			prop_name: xxx
func (p *defPropInputMap) UnmarshalYAML(value *yaml.Node) error {
	*p = make(defPropInputMap)

	// ScalarNode か MappingNode のみ許可
	if value.Kind == yaml.ScalarNode {
		(*p)["name"] = defPropValueInput{Value: value.Value}
	} else if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content); i += 2 {
			di := &defPropInput{}
			if err := value.Content[i].Decode(&di.Name); err != nil {
				return err
			}
			if err := value.Content[i+1].Decode(&di.Value); err != nil {
				return err
			}
			(*p)[di.Name] = di.Value
		}
	} else {
		return fmt.Errorf(errorInvalidKind+"\n"+errorInvalidDefs, value.Kind)
	}

	return nil
}

// defPropValueInput ... def の各プロパティの値
type defPropValueInput struct {
	IsSlice bool
	Value   string
	Values  []string
}

// def のプロパティの値の部分
func (p *defPropValueInput) UnmarshalYAML(value *yaml.Node) error {
	// ScalarNode か SequenceNode のみ許可
	if value.Kind == yaml.ScalarNode {
		p.Value = value.Value
	} else if value.Kind == yaml.SequenceNode {
		p.IsSlice = true
		p.Values = make([]string, len(value.Content))
		for i, v := range value.Content {
			if err := v.Decode(&p.Values[i]); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf(errorInvalidKind+"\n"+errorInvalidDefs, value.Kind)
	}
	return nil
}

// extendPropInput ... extends の各プロパティ
type extendPropInput struct {
	Name string
	Type string
}

// extendsInput ... extends を構造化したもの
type extendsInput struct {
	Name       string
	IsTemplate bool
	Props      []*extendPropInput
}

func (p *extendsInput) UnmarshalYAML(value *yaml.Node) error {
	// ScalarNode か MappingNode のみ許可
	if value.Kind == yaml.ScalarNode {
		*p = extendsInput{
			Name:       value.Value,
			IsTemplate: true,
			Props:      nil,
		}
		return nil
	}

	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`extends` must contain YAML mapping.")
	}

	*p = extendsInput{
		Name:       "",
		IsTemplate: false,
		Props:      make([]*extendPropInput, len(value.Content)/2),
	}

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

		p.Props[i/2] = &extendPropInput{
			Name: exName,
			Type: typeName,
		}
	}
	return nil
}

// extendsDefInput ... extends_defs を構造化したもの
type extendsDefInput = extendsInput

// extendsDefInputList ... extendsDefInput のリスト
type extendsDefInputList []*extendsDefInput

// extends_defs: の部分
func (p *extendsDefInputList) UnmarshalYAML(value *yaml.Node) error {
	// MappingNode のみ許可
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`extends_defs` must contain YAML mapping.")
	}

	*p = make([]*extendsDefInput, len(value.Content)/2)

	for i := 0; i < len(value.Content); i += 2 {
		ei := &extendsDefInput{}
		var keyName string

		// key
		if err := value.Content[i].Decode(&keyName); err != nil {
			return err
		}

		// value
		if err := value.Content[i+1].Decode(&ei); err != nil {
			return err
		}

		// extends_defs の中で key: scalar 指定があった場合エラー
		if ei.IsTemplate {
			return fmt.Errorf("`extends_defs` > `%v` must contain YAML mapping.", keyName)
		}

		ei.Name = keyName
		(*p)[i/2] = ei
	}
	return nil
}
