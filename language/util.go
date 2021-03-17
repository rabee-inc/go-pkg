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
		// 日本語
		return string(KeyJapanese)
	case strings.HasPrefix(lang, string(KeyEnglish)):
		// 英語
		return string(KeyEnglish)
	case strings.HasPrefix(lang, string(KeyTraditionalChinese)):
		// 中国語(繁)
		return string(KeyTraditionalChinese)
	case strings.HasPrefix(lang, string(keyTraditionalTaiwanChinese)):
		// 中国語(繁)の別パターン
		return string(KeyTraditionalChinese)
	default:
		// 不明の場合は日本語
		return string(KeyJapanese)
	}
}
