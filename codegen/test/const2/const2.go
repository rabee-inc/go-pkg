// Code generated by const2.yaml DO NOT EDIT.

package const2

const CheckSum = "27a51e9211da89fdc2a4717efc93bb5d94ade4d638815b77facf51bcda100f5f"

type ConstantMetaData[T comparable] struct {
	ID   T      `json:"id"`
	Name string `json:"name"`
}

type WithCategoryMetaDataProps struct {
	Name     string   `json:"name"`
	Category Category `json:"category"`
}

type WithCategoryMetaData[T WithCategory] struct {
	ID T `json:"id"`
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
	ID T `json:"id"`
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
	SkillMagicGuard Skill = "magic_guard"
	SkillBlackMagic Skill = "black_magic"
	SkillWhiteMagic Skill = "white_magic"
)

type SkillMetaData ConstantMetaData[Skill]

var Skills = []*SkillMetaData{
	{
		ID:   SkillMagicGuard,
		Name: "魔法無効",
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

// ExtendsDefsTestA ... test extends defs
type ExtendsDefsTestA string

func (c ExtendsDefsTestA) String() string {
	return string(c)
}

func (c ExtendsDefsTestA) Props() (*WithCategoryMetaDataProps, bool) {
	if m, ok := c.Meta(); ok {
		return m.WithCategoryMetaDataProps, ok
	}
	return nil, false
}

func (c ExtendsDefsTestA) Meta() (*ExtendsDefsTestAMetaData, bool) {
	m, ok := ExtendsDefsTestAMap[c]
	return m, ok
}

func (c ExtendsDefsTestA) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	ExtendsDefsTestAV1 ExtendsDefsTestA = "v1"
)

type ExtendsDefsTestAMetaData WithCategoryMetaData[ExtendsDefsTestA]

var ExtendsDefsTestAs = []*ExtendsDefsTestAMetaData{
	{
		ID: ExtendsDefsTestAV1,
		WithCategoryMetaDataProps: &WithCategoryMetaDataProps{
			Name:     "v1",
			Category: CategoryFood,
		},
	},
}

var ExtendsDefsTestAMap map[ExtendsDefsTestA]*ExtendsDefsTestAMetaData

// ExtendsDefsTestB ... test extends defs
type ExtendsDefsTestB string

func (c ExtendsDefsTestB) String() string {
	return string(c)
}

func (c ExtendsDefsTestB) Props() (*WithCategoryMetaDataProps, bool) {
	if m, ok := c.Meta(); ok {
		return m.WithCategoryMetaDataProps, ok
	}
	return nil, false
}

func (c ExtendsDefsTestB) Meta() (*ExtendsDefsTestBMetaData, bool) {
	m, ok := ExtendsDefsTestBMap[c]
	return m, ok
}

func (c ExtendsDefsTestB) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	ExtendsDefsTestBV1 ExtendsDefsTestB = "v1"
)

type ExtendsDefsTestBMetaData WithCategoryMetaData[ExtendsDefsTestB]

var ExtendsDefsTestBs = []*ExtendsDefsTestBMetaData{
	{
		ID: ExtendsDefsTestBV1,
		WithCategoryMetaDataProps: &WithCategoryMetaDataProps{
			Name:     "v1",
			Category: CategoryBag,
		},
	},
}

var ExtendsDefsTestBMap map[ExtendsDefsTestB]*ExtendsDefsTestBMetaData

// Player ... プレイヤー
type Player string

func (c Player) String() string {
	return string(c)
}

func (c Player) Props() (*CharacterStatusMetaDataProps, bool) {
	if m, ok := c.Meta(); ok {
		return m.CharacterStatusMetaDataProps, ok
	}
	return nil, false
}

func (c Player) Meta() (*PlayerMetaData, bool) {
	m, ok := PlayerMap[c]
	return m, ok
}

func (c Player) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	PlayerMagician Player = "magician"
)

type PlayerMetaData CharacterStatusMetaData[Player]

var Players = []*PlayerMetaData{
	{
		ID: PlayerMagician,
		CharacterStatusMetaDataProps: &CharacterStatusMetaDataProps{
			Name:   "黒魔法使い",
			Power:  1,
			Speed:  4,
			Detail: "魔法が使える",
			Skills: []Skill{SkillBlackMagic, SkillWhiteMagic},
		},
	},
}

var PlayerMap map[Player]*PlayerMetaData

// Enemy ... 敵
type Enemy string

func (c Enemy) String() string {
	return string(c)
}

func (c Enemy) Props() (*CharacterStatusMetaDataProps, bool) {
	if m, ok := c.Meta(); ok {
		return m.CharacterStatusMetaDataProps, ok
	}
	return nil, false
}

func (c Enemy) Meta() (*EnemyMetaData, bool) {
	m, ok := EnemyMap[c]
	return m, ok
}

func (c Enemy) Name() string {
	if m, ok := c.Meta(); ok {
		return m.Name
	}
	return ""
}

const (
	EnemyWolf   Enemy = "wolf"
	EnemyDragon Enemy = "dragon"
)

type EnemyMetaData CharacterStatusMetaData[Enemy]

var Enemies = []*EnemyMetaData{
	{
		ID: EnemyWolf,
		CharacterStatusMetaDataProps: &CharacterStatusMetaDataProps{
			Name:   "狼",
			Power:  10,
			Speed:  10,
			Detail: "凶暴",
			Skills: []Skill{},
		},
	},
	{
		ID: EnemyDragon,
		CharacterStatusMetaDataProps: &CharacterStatusMetaDataProps{
			Name:   "ドラゴン",
			Power:  10,
			Speed:  10,
			Detail: "凶暴",
			Skills: []Skill{SkillMagicGuard},
		},
	},
}

var EnemyMap map[Enemy]*EnemyMetaData

type Constants struct {
	Babies            []*BabyMetaData                                `json:"babies"`
	Baby              map[Baby]*BabyMetaData                         `json:"baby"`
	Toys              []*ToyMetaData                                 `json:"toys"`
	Toy               map[Toy]*ToyMetaData                           `json:"toy"`
	Oses              []*OsMetaData                                  `json:"oses"`
	Os                map[Os]*OsMetaData                             `json:"os"`
	Skills            []*SkillMetaData                               `json:"skills"`
	Skill             map[Skill]*SkillMetaData                       `json:"skill"`
	Categories        []*CategoryMetaData                            `json:"categories"`
	Category          map[Category]*CategoryMetaData                 `json:"category"`
	ExtendsDefsTestAs []*ExtendsDefsTestAMetaData                    `json:"extends_defs_test_as"`
	ExtendsDefsTestA  map[ExtendsDefsTestA]*ExtendsDefsTestAMetaData `json:"extends_defs_test_a"`
	ExtendsDefsTestBs []*ExtendsDefsTestBMetaData                    `json:"extends_defs_test_bs"`
	ExtendsDefsTestB  map[ExtendsDefsTestB]*ExtendsDefsTestBMetaData `json:"extends_defs_test_b"`
	Players           []*PlayerMetaData                              `json:"players"`
	Player            map[Player]*PlayerMetaData                     `json:"player"`
	Enemies           []*EnemyMetaData                               `json:"enemies"`
	Enemy             map[Enemy]*EnemyMetaData                       `json:"enemy"`
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

	extendsDefsTestA := []any{}
	for _, v := range c.ExtendsDefsTestAs {
		extendsDefsTestA = append(extendsDefsTestA, v.ID)
	}

	extendsDefsTestB := []any{}
	for _, v := range c.ExtendsDefsTestBs {
		extendsDefsTestB = append(extendsDefsTestB, v.ID)
	}

	player := []any{}
	for _, v := range c.Players {
		player = append(player, v.ID)
	}

	enemy := []any{}
	for _, v := range c.Enemies {
		enemy = append(enemy, v.ID)
	}

	return [][]any{
		baby,
		toy,
		os,
		skill,
		category,
		extendsDefsTestA,
		extendsDefsTestB,
		player,
		enemy,
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

	ExtendsDefsTestAMap = map[ExtendsDefsTestA]*ExtendsDefsTestAMetaData{}
	for _, v := range ExtendsDefsTestAs {
		ExtendsDefsTestAMap[v.ID] = v
	}

	ExtendsDefsTestBMap = map[ExtendsDefsTestB]*ExtendsDefsTestBMetaData{}
	for _, v := range ExtendsDefsTestBs {
		ExtendsDefsTestBMap[v.ID] = v
	}

	PlayerMap = map[Player]*PlayerMetaData{}
	for _, v := range Players {
		PlayerMap[v.ID] = v
	}

	EnemyMap = map[Enemy]*EnemyMetaData{}
	for _, v := range Enemies {
		EnemyMap[v.ID] = v
	}

	ConstantsData = &Constants{
		Babies:            Babies,
		Baby:              BabyMap,
		Toys:              Toys,
		Toy:               ToyMap,
		Oses:              Oses,
		Os:                OsMap,
		Skills:            Skills,
		Skill:             SkillMap,
		Categories:        Categories,
		Category:          CategoryMap,
		ExtendsDefsTestAs: ExtendsDefsTestAs,
		ExtendsDefsTestA:  ExtendsDefsTestAMap,
		ExtendsDefsTestBs: ExtendsDefsTestBs,
		ExtendsDefsTestB:  ExtendsDefsTestBMap,
		Players:           Players,
		Player:            PlayerMap,
		Enemies:           Enemies,
		Enemy:             EnemyMap,
	}
}
