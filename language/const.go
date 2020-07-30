package language

// Text ... 文言
type Text map[Key]string

// Key ... 言語の種類
type Key string

const (
	// KeyJapanese ... 言語: 日本語
	KeyJapanese Key = "ja"
	// KeyEnglish ... 言語: 英語
	KeyEnglish Key = "en"
	// KeyTraditionalChinese ... 言語: 中国語(繁)
	KeyTraditionalChinese Key = "zh-hant"
)
