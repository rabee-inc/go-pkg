package cloudfirestore

import (
	"reflect"

	"cloud.google.com/go/firestore"
)

func SetDocByDst(dst any, ref *firestore.DocumentRef) {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}
	SetDocByDsts(rv, rv.Elem().Type(), ref)
}

func SetDocByDsts(rv reflect.Value, rt reflect.Type, ref *firestore.DocumentRef) {
	if rt.Kind() != reflect.Struct {
		return
	}
	dst := rv.Elem()
	for i := range rt.NumField() {
		f := rt.Field(i)
		switch f.Tag.Get("cloudfirestore") {
		case "id":
			if f.Type.Kind() == reflect.String {
				dst.Field(i).SetString(ref.ID)
			}
		case "ref":
			if f.Type.Kind() == reflect.Pointer {
				dst.Field(i).Set(reflect.ValueOf(ref))
			}
		case "parent_id":
			if p := ancestorDoc(ref, 1); p != nil && f.Type.Kind() == reflect.String {
				dst.Field(i).SetString(p.ID)
			}
		case "parent_parent_id":
			if p := ancestorDoc(ref, 2); p != nil && f.Type.Kind() == reflect.String {
				dst.Field(i).SetString(p.ID)
			}
		}
	}
}

// ancestorDoc ... depth 階層上の親ドキュメントを取得する(存在しない場合は nil)
func ancestorDoc(ref *firestore.DocumentRef, depth int) *firestore.DocumentRef {
	doc := ref
	for range depth {
		if doc == nil || doc.Parent == nil {
			return nil
		}
		doc = doc.Parent.Parent
	}
	return doc
}

func SetEmptyBySlice(dst any) {
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

func SetEmptyBySlices(rv reflect.Value, rt reflect.Type) {
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

func SetEmptyByMap(dst any) {
	rv := reflect.Indirect(reflect.ValueOf(dst))
	rt := rv.Type()
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.Type.Kind() == reflect.Map && rv.Field(i).Len() == 0 {
				mp := reflect.MakeMap(f.Type)
				m := reflect.Indirect(mp)
				rv.Field(i).Set(m)
				continue
			}
		}
	}
}

func SetEmptyByMaps(rv reflect.Value, rt reflect.Type) {
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.Type.Kind() == reflect.Map && rv.Elem().Field(i).Len() == 0 {
				mp := reflect.MakeMap(f.Type)
				m := reflect.Indirect(mp)
				rv.Elem().Field(i).Set(m)
				continue
			}
		}
	}
}
