package util

// AssignIfNotNil ... src が nil でない場合に dest に代入する
func AssignIfNotNil[T any](dest *T, src *T) bool {
	if src != nil {
		*dest = *src
		return true
	}
	return false
}
