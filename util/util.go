package util

// AssignIfNotNil ... src が nil でない場合に dest に代入する
func AssignIfNotNil[T any](dest *T, src *T) {
	if src != nil {
		*dest = *src
	}
}
