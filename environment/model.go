package environment

// プロジェクト
type Project struct {
	Local      string `json:"local"`
	Staging    string `json:"staging"`
	Production string `json:"production"`
}

// 値
type Variable struct {
	Local      map[string]string `yaml:"local"`
	Staging    map[string]string `yaml:"staging"`
	Production map[string]string `yaml:"production"`
}
