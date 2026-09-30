package validation

// 値がゼロ値かどうかを判断する(true: ゼロ値, false: ゼロ値以外)
func IsZero(val any) bool {
	switch v := val.(type) {
	case nil:
		return true
	case int:
		return v == 0
	case int64:
		return v == 0
	case float64:
		return v == 0
	case string:
		return v == ""
	case bool:
		return !v
	default:
		return val == nil
	}
}
