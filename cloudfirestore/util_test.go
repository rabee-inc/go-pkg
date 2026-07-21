package cloudfirestore_test

import (
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/rabee-inc/go-pkg/cloudfirestore"
)

// createMockDocumentRef creates a mock DocumentRef with nested parents
func createMockDocumentRef() *firestore.DocumentRef {
	// Create a mock client (nil is fine for testing)
	var client *firestore.Client

	// Create nested collection/document structure:
	// collection/parent_parent_doc/collection/parent_doc/collection/doc
	col := client.Collection("root_collection")
	parentParentDoc := col.Doc("parent_parent_id")
	parentCol := parentParentDoc.Collection("parent_collection")
	parentDoc := parentCol.Doc("parent_id")
	targetCol := parentDoc.Collection("target_collection")
	doc := targetCol.Doc("test_id")

	return doc
}

func TestSetDocByDst(t *testing.T) {
	type testStruct struct {
		ID             string                     `cloudfirestore:"id"`
		Ref            *firestore.DocumentRef     `cloudfirestore:"ref"`
		ParentID       string                     `cloudfirestore:"parent_id"`
		ParentParentID string                     `cloudfirestore:"parent_parent_id"`
		Name           string                     // no tag
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ID != "test_id" {
					t.Errorf("ID = %v, want %v", dst.ID, "test_id")
				}
			},
		},
		{
			name: "ref タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Ref == nil {
					t.Error("Ref is nil")
				} else if dst.Ref.ID != "test_id" {
					t.Errorf("Ref.ID = %v, want %v", dst.Ref.ID, "test_id")
				}
			},
		},
		{
			name: "parent_id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ParentID != "parent_id" {
					t.Errorf("ParentID = %v, want %v", dst.ParentID, "parent_id")
				}
			},
		},
		{
			name: "parent_parent_id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ParentParentID != "parent_parent_id" {
					t.Errorf("ParentParentID = %v, want %v", dst.ParentParentID, "parent_parent_id")
				}
			},
		},
		{
			name: "タグなしフィールドは変更されない",
			dst:  &testStruct{Name: "original"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Name != "original" {
					t.Errorf("Name = %v, want %v", dst.Name, "original")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := createMockDocumentRef()
			cloudfirestore.SetDocByDst(tt.dst, ref)
			tt.validate(t, tt.dst)
		})
	}
}

func TestSetDocByDsts(t *testing.T) {
	type testStruct struct {
		ID             string                     `cloudfirestore:"id"`
		Ref            *firestore.DocumentRef     `cloudfirestore:"ref"`
		ParentID       string                     `cloudfirestore:"parent_id"`
		ParentParentID string                     `cloudfirestore:"parent_parent_id"`
		Name           string                     // no tag
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ID != "test_id" {
					t.Errorf("ID = %v, want %v", dst.ID, "test_id")
				}
			},
		},
		{
			name: "ref タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Ref == nil {
					t.Error("Ref is nil")
				} else if dst.Ref.ID != "test_id" {
					t.Errorf("Ref.ID = %v, want %v", dst.Ref.ID, "test_id")
				}
			},
		},
		{
			name: "parent_id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ParentID != "parent_id" {
					t.Errorf("ParentID = %v, want %v", dst.ParentID, "parent_id")
				}
			},
		},
		{
			name: "parent_parent_id タグが正しく設定される",
			dst:  &testStruct{Name: "test"},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.ParentParentID != "parent_parent_id" {
					t.Errorf("ParentParentID = %v, want %v", dst.ParentParentID, "parent_parent_id")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := createMockDocumentRef()
			rv := reflect.ValueOf(tt.dst)
			rt := rv.Elem().Type()
			cloudfirestore.SetDocByDsts(rv, rt, ref)
			tt.validate(t, tt.dst)
		})
	}
}

func TestSetEmptyBySlice(t *testing.T) {
	type testStruct struct {
		Name       string
		Slice      []string
		SliceWithValue []int
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "nil スライスが空スライスに初期化される",
			dst: &testStruct{
				Name:  "test",
				Slice: nil,
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Slice == nil {
					t.Error("Slice should not be nil")
				}
				if len(dst.Slice) != 0 {
					t.Errorf("Slice length = %v, want 0", len(dst.Slice))
				}
			},
		},
		{
			name: "空スライスはそのまま",
			dst: &testStruct{
				Name:  "test",
				Slice: []string{},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Slice == nil {
					t.Error("Slice should not be nil")
				}
				if len(dst.Slice) != 0 {
					t.Errorf("Slice length = %v, want 0", len(dst.Slice))
				}
			},
		},
		{
			name: "値があるスライスは変更されない",
			dst: &testStruct{
				Name:           "test",
				SliceWithValue: []int{1, 2, 3},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if len(dst.SliceWithValue) != 3 {
					t.Errorf("SliceWithValue length = %v, want 3", len(dst.SliceWithValue))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloudfirestore.SetEmptyBySlice(tt.dst)
			tt.validate(t, tt.dst)
		})
	}
}

func TestSetEmptyBySlices(t *testing.T) {
	type testStruct struct {
		Name  string
		Slice []string
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "nil スライスが空スライスに初期化される",
			dst: &testStruct{
				Name:  "test",
				Slice: nil,
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Slice == nil {
					t.Error("Slice should not be nil")
				}
				if len(dst.Slice) != 0 {
					t.Errorf("Slice length = %v, want 0", len(dst.Slice))
				}
			},
		},
		{
			name: "空スライスはそのまま",
			dst: &testStruct{
				Name:  "test",
				Slice: []string{},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Slice == nil {
					t.Error("Slice should not be nil")
				}
				if len(dst.Slice) != 0 {
					t.Errorf("Slice length = %v, want 0", len(dst.Slice))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := reflect.ValueOf(tt.dst)
			rt := rv.Elem().Type()
			cloudfirestore.SetEmptyBySlices(rv, rt)
			tt.validate(t, tt.dst)
		})
	}
}

func TestSetEmptyByMap(t *testing.T) {
	type testStruct struct {
		Name          string
		Map           map[string]string
		MapWithValue  map[string]int
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "nil マップが空マップに初期化される",
			dst: &testStruct{
				Name: "test",
				Map:  nil,
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Map == nil {
					t.Error("Map should not be nil")
				}
				if len(dst.Map) != 0 {
					t.Errorf("Map length = %v, want 0", len(dst.Map))
				}
			},
		},
		{
			name: "空マップはそのまま",
			dst: &testStruct{
				Name: "test",
				Map:  map[string]string{},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Map == nil {
					t.Error("Map should not be nil")
				}
				if len(dst.Map) != 0 {
					t.Errorf("Map length = %v, want 0", len(dst.Map))
				}
			},
		},
		{
			name: "値があるマップは変更されない",
			dst: &testStruct{
				Name:         "test",
				MapWithValue: map[string]int{"a": 1, "b": 2},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if len(dst.MapWithValue) != 2 {
					t.Errorf("MapWithValue length = %v, want 2", len(dst.MapWithValue))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloudfirestore.SetEmptyByMap(tt.dst)
			tt.validate(t, tt.dst)
		})
	}
}

func TestSetEmptyByMaps(t *testing.T) {
	type testStruct struct {
		Name string
		Map  map[string]string
	}

	tests := []struct {
		name     string
		dst      *testStruct
		validate func(*testing.T, *testStruct)
	}{
		{
			name: "nil マップが空マップに初期化される",
			dst: &testStruct{
				Name: "test",
				Map:  nil,
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Map == nil {
					t.Error("Map should not be nil")
				}
				if len(dst.Map) != 0 {
					t.Errorf("Map length = %v, want 0", len(dst.Map))
				}
			},
		},
		{
			name: "空マップはそのまま",
			dst: &testStruct{
				Name: "test",
				Map:  map[string]string{},
			},
			validate: func(t *testing.T, dst *testStruct) {
				if dst.Map == nil {
					t.Error("Map should not be nil")
				}
				if len(dst.Map) != 0 {
					t.Errorf("Map length = %v, want 0", len(dst.Map))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := reflect.ValueOf(tt.dst)
			rt := rv.Elem().Type()
			cloudfirestore.SetEmptyByMaps(rv, rt)
			tt.validate(t, tt.dst)
		})
	}
}
