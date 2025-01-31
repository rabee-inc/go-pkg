# 定数自動生成

## 生成処理の書き方

以下のコードは `./defs` ディレクトリ内の定義ファイルをすべて読み込み定数生成するコードです。

このコードは `main.go` という名前で配置し、そのディレクトリで

```bash
go generate
```

と実行することで、main関数が実行され定数が自動生成されます。

```go:main.go

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


```

## yaml の書き方

基本的にスネークケースで記載を推奨。キャメルケースには対応していません。

## 注意事項

同じディレクトリに同じpackage名で出力することはできません

### settings

```yaml
settings:
  package: パッケージ名
  output: 出力先 (go generate を実行するディレクトリからの相対パス)
```

### types

```yaml
types:
  定数の型名:
    comment: コメント # 必須
    only_backend: true # 任意: バックエンドでしか使用しない値かどうか。true の場合はフロントに返す値に含めません
    type: int # 任意 (int | int64 | string | float) (デフォルト = string)
    extends: # 任意: meta data に追加するプロパティ。 (ID と Name 以外に追加でプロパティを含める場合に使用してください)
      プロパティ名: 型名 # 必須: 型名 (int | int64 | string | float | またそれぞれのslice)
    defs: # 別途記載

```

### defs

以下の3種類の書き方があります。

#### ショートハンド

```yaml
変数名: Name の値
```

**入力例**

```yaml
types:
  comment: アイテム
  item:
    defs:
      private: 非公開
      public: 公開 
```

**出力**

```go
// Item ... アイテム
type Item string

const (
	ItemPrivate Item = "private"
	ItemPublic  Item = "public"
)

type ItemMetaData ConstantMetaData[Item]

var Items = []*ItemMetaData{
	{
		ID:   ItemPrivate,
		Name: "非公開",
	},
	{
		ID:   ItemPublic,
		Name: "公開",
	},
}

var ItemMap map[Item]*ItemMetaData

```

#### ID を変更する場合

```yaml
変数名:
  id: ID の値
  name: Name の値
```


**入力例**

```yaml
types:
  comment: アイテム
  type: int
  item:
    defs:
      private:
        id: "1" # (数値を入力する場合は "" で囲んでください)
        name: 非公開
      public:
        id: "2"
        name: 公開 
```

**出力**

```go
// Item ... アイテム
type Item int

const (
	ItemPrivate Item = 1
	ItemPublic  Item = 2
)

type ItemMetaData ConstantMetaData[Item]

var Items = []*ItemMetaData{
	{
		ID:   ItemPrivate,
		Name: "非公開",
	},
	{
		ID:   ItemPublic,
		Name: "公開",
	},
}

var ItemMap map[Item]*ItemMetaData

```

#### ID と Name 以外にプロパティを追加する場合

`extends` を定義してその値をプロパティに追加します。

**入力例**

```yaml
types:
  comment: アイテム
  extends: 
    color: string
  item:
    defs:
      private:
        name: 非公開
        color: red
      public:
        name: 公開 
        color: green
```

**出力**

```go
// Item ... アイテム
type Item string

const (
	ItemPrivate Item = "private"
	ItemPublic  Item = "public"
)

type ItemMetaData ConstantMetaData[Item]

var Items = []*ItemMetaData{
	{
		ID:   ItemPrivate,
		Name: "非公開",
		Color: "red"
	},
	{
		ID:   ItemPublic,
		Name: "公開",
		Color: "green"
	},
}

var ItemMap map[Item]*ItemMetaData

```
