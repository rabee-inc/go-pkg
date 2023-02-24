package cloudfirestore

import (
	"reflect"

	"cloud.google.com/go/firestore"
)

func setDocByDst(dst any, ref *firestore.DocumentRef) {
	rv := reflect.Indirect(reflect.ValueOf(dst))
	rt := rv.Type()
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			tag := f.Tag.Get("cloudfirestore")
			if tag == "id" && f.Type.Kind() == reflect.String {
				rv.Field(i).SetString(ref.ID)
				continue
			}
			if tag == "ref" && f.Type.Kind() == reflect.Ptr {
				rv.Field(i).Set(reflect.ValueOf(ref))
				continue
			}
		}
	}
}

func setDocByDsts(rv reflect.Value, rt reflect.Type, ref *firestore.DocumentRef) {
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			tag := f.Tag.Get("cloudfirestore")
			if tag == "id" && f.Type.Kind() == reflect.String {
				rv.Elem().Field(i).SetString(ref.ID)
				continue
			}
			if tag == "ref" && f.Type.Kind() == reflect.Ptr {
				rv.Elem().Field(i).Set(reflect.ValueOf(ref))
				continue
			}
		}
	}
}

func setEmptyBySlice(dst any) {
	rv := reflect.Indirect(reflect.ValueOf(dst))
	rt := rv.Type()
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.Type.Kind() == reflect.Slice && rv.Field(i).Len() == 0 {
				sp := reflect.MakeSlice(f.Type, 0, 0)
				s := reflect.Indirect(sp)
				rv.Field(i).Set(s)
				continue
			}
		}
	}
}

func setEmptyBySlices(rv reflect.Value, rt reflect.Type) {
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.Type.Kind() == reflect.Slice && rv.Elem().Field(i).Len() == 0 {
				sp := reflect.MakeSlice(f.Type, 0, 0)
				s := reflect.Indirect(sp)
				rv.Elem().Field(i).Set(s)
				continue
			}
		}
	}
}
