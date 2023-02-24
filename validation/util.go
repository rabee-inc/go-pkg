package validation

// 値がゼロ値かどうかを判断する(true: ゼロ値, false: ゼロ値以外)
func IsZero(val any) bool {
	switch val.(type) {
	case nil:
		return true
	case int:
		if val.(int) == 0 {
			return true
		}
		return false
	case int64:
		if val.(int64) == 0 {
			return true
		}
		return false
	case float64:
		if val.(float64) == 0 {
			return true
		}
		return false
	case string:
		if val.(string) == "" {
			return true
		}
		return false
	case bool:
		if !val.(bool) {
			return true
		}
		return false
	default:
		if val == nil {
			return true
		}
		return false
	}
}
