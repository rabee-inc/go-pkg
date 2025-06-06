package codegen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rabee-inc/go-pkg/codegen"
	"github.com/rabee-inc/go-pkg/codegen/test/const1"
	"github.com/rabee-inc/go-pkg/codegen/test/const2"
)

func TestConst(t *testing.T) {
	// コンパイルが通るかというテスト
	meta, ok := const1.AnimalCat.Meta()
	fmt.Println(const1.AnimalCat.String())
	if ok {
		fmt.Println(meta.Name)
	}
	fmt.Println(const2.ColorRed.Name())

	for k, v := range const1.ExtendsTests {
		fmt.Println(k, v)
		fmt.Println(v.ID)
		fmt.Println(v.Name)
		fmt.Println(v.IntValue)
		fmt.Println(v.IntSliceValue)
		fmt.Println(v.Int64Value)
		fmt.Println(v.Int64SliceValue)
		fmt.Println(v.FloatValue)
		fmt.Println(v.FloatSliceValue)
		fmt.Println(v.StringValue)
		fmt.Println(v.StringSliceValue)
	}

	// other type
	{
		meta, _ := const1.TypeTestV1.Meta()
		meta2, _ := meta.ExtendsTest.Meta()
		fmt.Println(meta2.IntValue)
	}

	// other type slice
	{
		meta, _ := const1.TypeTestV1.Meta()
		for _, animal := range meta.Animals {
			animalMeta, _ := animal.Meta()
			fmt.Println(animalMeta.Name)
		}
	}

	// extends_defs Props
	{
		props, _ := const2.ExtendsDefsTestAV1.Props()
		fmt.Println("Props", props)

		withCategories := []const2.WithCategory{
			const2.ExtendsDefsTestAV1,
			const2.ExtendsDefsTestBV1,
		}

		for _, v := range withCategories {
			props, _ := v.Props()
			fmt.Println("Props", props)
		}

		meta, _ := const2.PlayerMagician.Meta()
		fmt.Println("Meta", meta, meta.CharacterStatusMetaDataProps)

		characters := []const2.CharacterStatus{
			const2.PlayerMagician,
			const2.EnemyDragon,
			const2.EnemyWolf,
		}
		for _, v := range characters {
			props, _ := v.Props()
			fmt.Println("Props", props)
			skillNames := []string{}
			for _, skill := range props.Skills {
				skillNames = append(skillNames, skill.Name())
			}
			fmt.Println("SkillNames", skillNames)
		}
	}
}

func TestGenerate(t *testing.T) {
	defsDir := "./test/defs"
	// defsDir 内のすべての yaml ファイルを読み込む
	files, err := os.ReadDir(defsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(defsDir, file.Name())
			outputPath := codegen.ExportByYaml(path)
			// delete file
			os.Remove(outputPath)
		}
	}
}
