package util

import "reflect"

// AssignIfNotNil ... src が nil でない場合に dest に代入する。代入前の値を含むpointerを返す
func AssignIfNotNil[T any](dest *T, src *T) *T {
	if src != nil {
		old := *dest
		*dest = *src
		return &old
	}
	return dest
}

// EachTaggedFields ... 指定のタグをがついてるフィールドをループする
func EachTaggedFields(param any, tagName string, callback func(tagValue string, reflectParam reflect.Value, fieldNum int) error) error {
	val := reflect.Indirect(reflect.ValueOf(param))
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)

		tag := typeField.Tag.Get(tagName)
		if tag == "" {
			continue
		}

		if err := callback(tag, val, i); err != nil {
			return err
		}
	}
	return nil
}
