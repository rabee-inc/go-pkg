package cloudfirestore_test

import (
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/rabee-inc/go-pkg/cloudfirestore"
)

// newColRef ... Client を介さずに CollectionRef を組み立てる
func newColRef(parent *firestore.DocumentRef, id string) *firestore.CollectionRef {
	path := "projects/test-project/databases/(default)/documents/" + id
	if parent != nil {
		path = parent.Path + "/" + id
	}
	return &firestore.CollectionRef{
		Parent: parent,
		ID:     id,
		Path:   path,
	}
}

// newDocRef ... Client を介さずに DocumentRef を組み立てる
func newDocRef(parent *firestore.CollectionRef, id string) *firestore.DocumentRef {
	return &firestore.DocumentRef{
		Parent: parent,
		ID:     id,
		Path:   parent.Path + "/" + id,
	}
}

// docRefFromPath ... Client を介さずに "col/doc/col/doc..." 形式のパスから DocumentRef を組み立てる
// 不正なパス(要素数が奇数、空の要素を含む)の場合は firestore.Client.Doc と同じく nil を返す
func docRefFromPath(path string) *firestore.DocumentRef {
	parts := strings.Split(path, "/")
	if len(parts)%2 != 0 {
		return nil
	}
	var doc *firestore.DocumentRef
	for i := 0; i < len(parts); i += 2 {
		if parts[i] == "" || parts[i+1] == "" {
			return nil
		}
		doc = newDocRef(newColRef(doc, parts[i]), parts[i+1])
	}
	return doc
}

// newNestedDocRef ... root_collection/parent_parent_id/parent_collection/parent_id/target_collection/test_id
func newNestedDocRef() *firestore.DocumentRef {
	parentParentDoc := newDocRef(newColRef(nil, "root_collection"), "parent_parent_id")
	parentDoc := newDocRef(newColRef(parentParentDoc, "parent_collection"), "parent_id")
	return newDocRef(newColRef(parentDoc, "target_collection"), "test_id")
}

type docTagStruct struct {
	ID             string                 `cloudfirestore:"id"`
	Ref            *firestore.DocumentRef `cloudfirestore:"ref"`
	ParentID       string                 `cloudfirestore:"parent_id"`
	ParentParentID string                 `cloudfirestore:"parent_parent_id"`
	Name           string                 // タグなし
}

func Test_Cloudfirestore_SetDocByDst(t *testing.T) {
	type args struct {
		dst any
		ref *firestore.DocumentRef
	}
	type want struct {
		id             string
		refID          string
		parentID       string
		parentParentID string
		name           string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 全ての階層のタグが設定される",
			args: args{
				dst: &docTagStruct{Name: "original"},
				ref: newNestedDocRef(),
			},
			want: want{
				id:             "test_id",
				refID:          "test_id",
				parentID:       "parent_id",
				parentParentID: "parent_parent_id",
				name:           "original",
			},
		},
		{
			name: "正常系: 親が浅い場合は親のIDは設定されない",
			args: args{
				dst: &docTagStruct{Name: "original"},
				ref: newDocRef(newColRef(nil, "root_collection"), "test_id"),
			},
			want: want{
				id:             "test_id",
				refID:          "test_id",
				parentID:       "",
				parentParentID: "",
				name:           "original",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cloudfirestore.SetDocByDst(tc.args.dst, tc.args.ref)

			dst := tc.args.dst.(*docTagStruct)
			if dst.ID != tc.want.id {
				t.Errorf("ID = %v, want %v", dst.ID, tc.want.id)
			}
			if dst.Ref == nil {
				t.Errorf("Ref = nil, want %v", tc.want.refID)
			} else if dst.Ref.ID != tc.want.refID {
				t.Errorf("Ref.ID = %v, want %v", dst.Ref.ID, tc.want.refID)
			}
			if dst.ParentID != tc.want.parentID {
				t.Errorf("ParentID = %v, want %v", dst.ParentID, tc.want.parentID)
			}
			if dst.ParentParentID != tc.want.parentParentID {
				t.Errorf("ParentParentID = %v, want %v", dst.ParentParentID, tc.want.parentParentID)
			}
			if dst.Name != tc.want.name {
				t.Errorf("Name = %v, want %v", dst.Name, tc.want.name)
			}
		})
	}
}

func Test_Cloudfirestore_SetDocByDst_Invalid(t *testing.T) {
	type args struct {
		dst any
	}
	type testCase struct {
		name string
		args args
	}

	tcs := []testCase{
		{
			name: "異常系: ポインタ以外を渡してもpanicしない",
			args: args{dst: docTagStruct{}},
		},
		{
			name: "異常系: nilポインタを渡してもpanicしない",
			args: args{dst: (*docTagStruct)(nil)},
		},
		{
			name: "異常系: 構造体以外のポインタを渡してもpanicしない",
			args: args{dst: new(string)},
		},
		{
			name: "異常系: nilを渡してもpanicしない",
			args: args{dst: nil},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("panic = %v, want no panic", r)
				}
			}()
			cloudfirestore.SetDocByDst(tc.args.dst, newNestedDocRef())
		})
	}
}

func Test_Cloudfirestore_SetDocByDsts(t *testing.T) {
	type args struct {
		dst *docTagStruct
		ref *firestore.DocumentRef
	}
	type want struct {
		id             string
		parentID       string
		parentParentID string
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: 全ての階層のタグが設定される",
			args: args{
				dst: &docTagStruct{},
				ref: newNestedDocRef(),
			},
			want: want{
				id:             "test_id",
				parentID:       "parent_id",
				parentParentID: "parent_parent_id",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.ValueOf(tc.args.dst)
			cloudfirestore.SetDocByDsts(rv, rv.Elem().Type(), tc.args.ref)

			if tc.args.dst.ID != tc.want.id {
				t.Errorf("ID = %v, want %v", tc.args.dst.ID, tc.want.id)
			}
			if tc.args.dst.ParentID != tc.want.parentID {
				t.Errorf("ParentID = %v, want %v", tc.args.dst.ParentID, tc.want.parentID)
			}
			if tc.args.dst.ParentParentID != tc.want.parentParentID {
				t.Errorf("ParentParentID = %v, want %v", tc.args.dst.ParentParentID, tc.want.parentParentID)
			}
		})
	}
}

type emptyTargetStruct struct {
	Name   string
	Slice  []string
	Ints   []int
	Map    map[string]string
	IntMap map[string]int
}

func Test_Cloudfirestore_SetEmptyBySlice(t *testing.T) {
	type args struct {
		dst *emptyTargetStruct
	}
	type want struct {
		sliceIsNil bool
		sliceLen   int
		intsIsNil  bool
		intsLen    int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: nilスライスが空スライスに初期化される",
			args: args{dst: &emptyTargetStruct{Name: "test", Slice: nil}},
			want: want{sliceIsNil: false, sliceLen: 0, intsIsNil: false, intsLen: 0},
		},
		{
			name: "正常系: 空スライスはそのまま",
			args: args{dst: &emptyTargetStruct{Name: "test", Slice: []string{}}},
			want: want{sliceIsNil: false, sliceLen: 0, intsIsNil: false, intsLen: 0},
		},
		{
			name: "正常系: 値があるスライスは変更されない",
			args: args{dst: &emptyTargetStruct{Name: "test", Ints: []int{1, 2, 3}}},
			want: want{sliceIsNil: false, sliceLen: 0, intsIsNil: false, intsLen: 3},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cloudfirestore.SetEmptyBySlice(tc.args.dst)

			if (tc.args.dst.Slice == nil) != tc.want.sliceIsNil {
				t.Errorf("Slice == nil = %v, want %v", tc.args.dst.Slice == nil, tc.want.sliceIsNil)
			}
			if len(tc.args.dst.Slice) != tc.want.sliceLen {
				t.Errorf("len(Slice) = %v, want %v", len(tc.args.dst.Slice), tc.want.sliceLen)
			}
			if (tc.args.dst.Ints == nil) != tc.want.intsIsNil {
				t.Errorf("Ints == nil = %v, want %v", tc.args.dst.Ints == nil, tc.want.intsIsNil)
			}
			if len(tc.args.dst.Ints) != tc.want.intsLen {
				t.Errorf("len(Ints) = %v, want %v", len(tc.args.dst.Ints), tc.want.intsLen)
			}
		})
	}
}

func Test_Cloudfirestore_SetEmptyBySlices(t *testing.T) {
	type args struct {
		dst *emptyTargetStruct
	}
	type want struct {
		sliceIsNil bool
		sliceLen   int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: nilスライスが空スライスに初期化される",
			args: args{dst: &emptyTargetStruct{Name: "test", Slice: nil}},
			want: want{sliceIsNil: false, sliceLen: 0},
		},
		{
			name: "正常系: 空スライスはそのまま",
			args: args{dst: &emptyTargetStruct{Name: "test", Slice: []string{}}},
			want: want{sliceIsNil: false, sliceLen: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.ValueOf(tc.args.dst)
			cloudfirestore.SetEmptyBySlices(rv, rv.Elem().Type())

			if (tc.args.dst.Slice == nil) != tc.want.sliceIsNil {
				t.Errorf("Slice == nil = %v, want %v", tc.args.dst.Slice == nil, tc.want.sliceIsNil)
			}
			if len(tc.args.dst.Slice) != tc.want.sliceLen {
				t.Errorf("len(Slice) = %v, want %v", len(tc.args.dst.Slice), tc.want.sliceLen)
			}
		})
	}
}

func Test_Cloudfirestore_SetEmptyByMap(t *testing.T) {
	type args struct {
		dst *emptyTargetStruct
	}
	type want struct {
		mapIsNil    bool
		mapLen      int
		intMapIsNil bool
		intMapLen   int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: nilマップが空マップに初期化される",
			args: args{dst: &emptyTargetStruct{Name: "test", Map: nil}},
			want: want{mapIsNil: false, mapLen: 0, intMapIsNil: false, intMapLen: 0},
		},
		{
			name: "正常系: 空マップはそのまま",
			args: args{dst: &emptyTargetStruct{Name: "test", Map: map[string]string{}}},
			want: want{mapIsNil: false, mapLen: 0, intMapIsNil: false, intMapLen: 0},
		},
		{
			name: "正常系: 値があるマップは変更されない",
			args: args{dst: &emptyTargetStruct{Name: "test", IntMap: map[string]int{"a": 1, "b": 2}}},
			want: want{mapIsNil: false, mapLen: 0, intMapIsNil: false, intMapLen: 2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cloudfirestore.SetEmptyByMap(tc.args.dst)

			if (tc.args.dst.Map == nil) != tc.want.mapIsNil {
				t.Errorf("Map == nil = %v, want %v", tc.args.dst.Map == nil, tc.want.mapIsNil)
			}
			if len(tc.args.dst.Map) != tc.want.mapLen {
				t.Errorf("len(Map) = %v, want %v", len(tc.args.dst.Map), tc.want.mapLen)
			}
			if (tc.args.dst.IntMap == nil) != tc.want.intMapIsNil {
				t.Errorf("IntMap == nil = %v, want %v", tc.args.dst.IntMap == nil, tc.want.intMapIsNil)
			}
			if len(tc.args.dst.IntMap) != tc.want.intMapLen {
				t.Errorf("len(IntMap) = %v, want %v", len(tc.args.dst.IntMap), tc.want.intMapLen)
			}
		})
	}
}

func Test_Cloudfirestore_SetEmptyByMaps(t *testing.T) {
	type args struct {
		dst *emptyTargetStruct
	}
	type want struct {
		mapIsNil bool
		mapLen   int
	}
	type testCase struct {
		name string
		args args
		want want
	}

	tcs := []testCase{
		{
			name: "正常系: nilマップが空マップに初期化される",
			args: args{dst: &emptyTargetStruct{Name: "test", Map: nil}},
			want: want{mapIsNil: false, mapLen: 0},
		},
		{
			name: "正常系: 空マップはそのまま",
			args: args{dst: &emptyTargetStruct{Name: "test", Map: map[string]string{}}},
			want: want{mapIsNil: false, mapLen: 0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rv := reflect.ValueOf(tc.args.dst)
			cloudfirestore.SetEmptyByMaps(rv, rv.Elem().Type())

			if (tc.args.dst.Map == nil) != tc.want.mapIsNil {
				t.Errorf("Map == nil = %v, want %v", tc.args.dst.Map == nil, tc.want.mapIsNil)
			}
			if len(tc.args.dst.Map) != tc.want.mapLen {
				t.Errorf("len(Map) = %v, want %v", len(tc.args.dst.Map), tc.want.mapLen)
			}
		})
	}
}
