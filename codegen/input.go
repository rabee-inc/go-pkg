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
	ExtendsDefs map[string]extendsInputList `yaml:"extends_defs"`
}

// typeInput ... types のそれぞれを構造化したもの
type typeInput struct {
	Name                string
	Comment             string           `yaml:"comment" validate:"required"`
	Type                string           `yaml:"type"`
	OnlyBackend         bool             `yaml:"only_backend"`
	Extends             extendsInputList `yaml:"extends"`
	ExtendsTemplateName string
	Defs                defInputList `yaml:"defs" validate:"required"`
}

// typeInputList ... types を構造化したもの
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

// defInput ... defs のそれぞれを構造化したもの
type defInput struct {
	Name       string
	OtherProps any
}

// defInputList ... defs を構造化したもの
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

// extendPropInput ... extends の各プロパティ
type extendPropInput struct {
	Name string
	Type string
}

// extendsInputList ... extends を構造化したもの
type extendsInputList []*extendPropInput

func (p *extendsInputList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("`types` must contain YAML mapping, has %v", value.Kind)
	}

	*p = make([]*extendPropInput, len(value.Content)/2)

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

		(*p)[i/2] = &extendPropInput{
			Name: exName,
			Type: typeName,
		}
	}
	return nil
}
