//go:generate go run .

package main

import (
	"os"
	"path/filepath"

	"github.com/rabee-inc/go-pkg/codegen"
)

const defsDir = "./defs"

func main() {
	// defsDir 内のすべての yaml ファイルを読み込む
	files, err := os.ReadDir(defsDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(defsDir, file.Name())
			codegen.ExportByYaml(path)
		}
	}
}
