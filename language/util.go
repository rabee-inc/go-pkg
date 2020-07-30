package language

import (
	"strings"
)

func trim(lang string) string {
	// 小文字にする
	lang = strings.ToLower(lang)

	// 先頭文字で判定する
	switch {
	case strings.HasPrefix(lang, string(KeyJapanese)):
		return string(KeyJapanese)
	case strings.HasPrefix(lang, string(KeyEnglish)):
		return string(KeyEnglish)
	case strings.HasPrefix(lang, string(KeyTraditionalChinese)):
		return string(KeyTraditionalChinese)
	default:
		return string(KeyJapanese)
	}
}
