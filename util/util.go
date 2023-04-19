package util

// AssignIfNotNil ... src が nil でない場合に dest に代入する。代入前の値を含むpointerを返す
func AssignIfNotNil[T any](dest *T, src *T) *T {
	if src != nil {
		old := *dest
		*dest = *src
		return &old
	}
	return dest
}
