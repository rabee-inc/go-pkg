package sliceutil

import "reflect"

// Stream ... スライス操作
type Stream struct {
	slice reflect.Value
}

/*
dst := StreamOf(hoges).
	Filter(func(hoge *Hoge) bool {
		return hoge.Num > 3
	}).Out().([]*Hoge)
*/
// Filter ... 要素のフィルタリング
func (s *Stream) Filter(fn interface{}) *Stream {
	frv := reflect.ValueOf(fn)
	srv := reflect.MakeSlice(s.slice.Type(), 0, 0)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		if out[0].Interface().(bool) {
			srv = reflect.Append(srv, rv)
		}
	}
	s.slice = srv
	return s
}

/*
dst := StreamOf(hoges).
	Map(func(hoge *Hoge) string {
		return hoge.ID
	}).Out().([]string)
*/
// Map ... 要素の変換
func (s *Stream) Map(fn interface{}) *Stream {
	frv := reflect.ValueOf(fn)
	srt := reflect.SliceOf(frv.Type().Out(0))
	srv := reflect.MakeSlice(srt, 0, 0)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		srv = reflect.Append(srv, out[0])
	}
	s.slice = srv
	return s
}

/*
dst := StreamOf(hoges).
	Reduce(func(dst int, num int) int {
		return dst + num
	}).(int)
*/
// Reduce ... 要素の集計
func (s *Stream) Reduce(fn interface{}) interface{} {
	frv := reflect.ValueOf(fn)
	rt := frv.Type().Out(0)
	dst := reflect.New(rt).Elem()
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{dst, rv})
		dst = out[0]
	}
	return dst.Interface()
}

/*
dst := StreamOf(hoges).
    Contains(func(hoge *Hoge) bool {
		return hoge.ID == "abc"
	})
*/
// Contains ... 要素の存在確認
func (s *Stream) Contains(fn interface{}) bool {
	frv := reflect.ValueOf(fn)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		if out[0].Interface().(bool) {
			return true
		}
	}
	return false
}

/*
dst := StreamOf(hoges).
    ForEach(func(hoge *Hoge) {
		hoge.ID = "abc"
	})
*/
// ForEach ... 要素のループ
func (s *Stream) ForEach(fn interface{}) *Stream {
	frv := reflect.ValueOf(fn)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		_ = frv.Call([]reflect.Value{rv})
	}
	return s
}

/*
dst := StreamOf(hoges).Count()
*/
// Count ... 要素数を取得
func (s *Stream) Count() int {
	return s.slice.Len()
}

// Out ... 結果を出力する
func (s *Stream) Out() interface{} {
	return s.slice.Interface()
}

// StreamOf ... スライスからスライス操作を作成する
func StreamOf(slice interface{}) *Stream {
	rv := reflect.ValueOf(slice)
	return &Stream{
		slice: rv,
	}
}
