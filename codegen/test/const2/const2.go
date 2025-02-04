// Code generated by const2.yaml DO NOT EDIT.

package const2

const CheckSum = "845e642b0f38c1f7d116ca848a6c1be30fb814f321aa20008ef946108c092880"

type ConstantMetaData[T comparable] struct {
	ID   T      `json:"id"`
	Name string `json:"name"`
}

type WithCategoryMetaDataProps struct {
	Name     string   `json:"name"`
	Category Category `json:"category"`
}

type WithCategoryMetaData[T WithCategory] struct {
	ID T
	*WithCategoryMetaDataProps
}

type WithCategory interface {
	Props() (*WithCategoryMetaDataProps, bool)
}

type CharacterStatusMetaDataProps struct {
	Name   string  `json:"name"`
	Power  int     `json:"power"`
	Speed  int     `json:"speed"`
	Detail string  `json:"detail"`
	Skills []Skill `json:"skills"`
}

type CharacterStatusMetaData[T CharacterStatus] struct {
	ID T
	*CharacterStatusMetaDataProps
}

type CharacterStatus interface {
	Props() (*CharacterStatusMetaDataProps, bool)
}

// Color ... 色
type Color int

func (c Color) Meta() (*ColorMetaData, bool) {
	m, ok := ColorMap[c]
	return m, ok
}

func (c Color) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	ColorRed  Color = 1
	ColorBlue Color = 2
)

type ColorMetaData ConstantMetaData[Color]

var Colors = []*ColorMetaData{
	{
		ID:   ColorRed,
		Name: "赤",
	},
	{
		ID:   ColorBlue,
		Name: "青",
	},
}

var ColorMap map[Color]*ColorMetaData

// Baby ... test: -y to -ies
type Baby string

func (c Baby) String() string {
	return string(c)
}

func (c Baby) Meta() (*BabyMetaData, bool) {
	m, ok := BabyMap[c]
	return m, ok
}

func (c Baby) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	BabyV1 Baby = "v1"
)

type BabyMetaData ConstantMetaData[Baby]

var Babies = []*BabyMetaData{
	{
		ID:   BabyV1,
		Name: "v1",
	},
}

var BabyMap map[Baby]*BabyMetaData

// Toy ... test: -y to -ys
type Toy string

func (c Toy) String() string {
	return string(c)
}

func (c Toy) Meta() (*ToyMetaData, bool) {
	m, ok := ToyMap[c]
	return m, ok
}

func (c Toy) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	ToyV1 Toy = "v1"
)

type ToyMetaData ConstantMetaData[Toy]

var Toys = []*ToyMetaData{
	{
		ID:   ToyV1,
		Name: "v1",
	},
}

var ToyMap map[Toy]*ToyMetaData

// Os ... test: -s to -es
type Os string

func (c Os) String() string {
	return string(c)
}

func (c Os) Meta() (*OsMetaData, bool) {
	m, ok := OsMap[c]
	return m, ok
}

func (c Os) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	OsV1 Os = "v1"
)

type OsMetaData ConstantMetaData[Os]

var Oses = []*OsMetaData{
	{
		ID:   OsV1,
		Name: "v1",
	},
}

var OsMap map[Os]*OsMetaData

// Skill ... スキル
type Skill string

func (c Skill) String() string {
	return string(c)
}

func (c Skill) Meta() (*SkillMetaData, bool) {
	m, ok := SkillMap[c]
	return m, ok
}

func (c Skill) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	SkillNoSkill    Skill = "no_skill"
	SkillBlackMagic Skill = "black_magic"
	SkillWhiteMagic Skill = "white_magic"
)

type SkillMetaData ConstantMetaData[Skill]

var Skills = []*SkillMetaData{
	{
		ID:   SkillNoSkill,
		Name: "スキルなし",
	},
	{
		ID:   SkillBlackMagic,
		Name: "黒魔法",
	},
	{
		ID:   SkillWhiteMagic,
		Name: "白魔法",
	},
}

var SkillMap map[Skill]*SkillMetaData

// Category ... カテゴリ
type Category string

func (c Category) String() string {
	return string(c)
}

func (c Category) Meta() (*CategoryMetaData, bool) {
	m, ok := CategoryMap[c]
	return m, ok
}

func (c Category) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	CategoryFood  Category = "food"
	CategoryShoes Category = "shoes"
	CategoryBag   Category = "bag"
)

type CategoryMetaData ConstantMetaData[Category]

var Categories = []*CategoryMetaData{
	{
		ID:   CategoryFood,
		Name: "フード",
	},
	{
		ID:   CategoryShoes,
		Name: "靴",
	},
	{
		ID:   CategoryBag,
		Name: "バッグ",
	},
}

var CategoryMap map[Category]*CategoryMetaData

type Constants struct {
	Babies     []*BabyMetaData                `json:"babies"`
	Baby       map[Baby]*BabyMetaData         `json:"baby"`
	Toys       []*ToyMetaData                 `json:"toys"`
	Toy        map[Toy]*ToyMetaData           `json:"toy"`
	Oses       []*OsMetaData                  `json:"oses"`
	Os         map[Os]*OsMetaData             `json:"os"`
	Skills     []*SkillMetaData               `json:"skills"`
	Skill      map[Skill]*SkillMetaData       `json:"skill"`
	Categories []*CategoryMetaData            `json:"categories"`
	Category   map[Category]*CategoryMetaData `json:"category"`
}

var ConstantsData *Constants

// deprecated use ConstIDs
func (c *Constants) GetConstIDs() [][]any {

	baby := []any{}
	for _, v := range c.Babies {
		baby = append(baby, v.ID)
	}

	toy := []any{}
	for _, v := range c.Toys {
		toy = append(toy, v.ID)
	}

	os := []any{}
	for _, v := range c.Oses {
		os = append(os, v.ID)
	}

	skill := []any{}
	for _, v := range c.Skills {
		skill = append(skill, v.ID)
	}

	category := []any{}
	for _, v := range c.Categories {
		category = append(category, v.ID)
	}

	return [][]any{
		baby,
		toy,
		os,
		skill,
		category,
	}
}

func (c *Constants) ConstIDs() [][]any {
	return c.GetConstIDs()
}

func init() {

	BabyMap = map[Baby]*BabyMetaData{}
	for _, v := range Babies {
		BabyMap[v.ID] = v
	}

	ToyMap = map[Toy]*ToyMetaData{}
	for _, v := range Toys {
		ToyMap[v.ID] = v
	}

	OsMap = map[Os]*OsMetaData{}
	for _, v := range Oses {
		OsMap[v.ID] = v
	}

	SkillMap = map[Skill]*SkillMetaData{}
	for _, v := range Skills {
		SkillMap[v.ID] = v
	}

	CategoryMap = map[Category]*CategoryMetaData{}
	for _, v := range Categories {
		CategoryMap[v.ID] = v
	}

	ConstantsData = &Constants{
		Babies:     Babies,
		Baby:       BabyMap,
		Toys:       Toys,
		Toy:        ToyMap,
		Oses:       Oses,
		Os:         OsMap,
		Skills:     Skills,
		Skill:      SkillMap,
		Categories: Categories,
		Category:   CategoryMap,
	}
}
