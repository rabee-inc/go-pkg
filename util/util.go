package util

func AssignIfNotNil[T any](dest *T, src *T) {
	if src != nil {
		*dest = *src
	}
}
