package cloudfirestore

import (
	"context"

	"github.com/rabee-inc/go-pkg/sliceutil"
	"golang.org/x/exp/maps"
)

type convertibleBatchGetterItem[D any] struct {
	ID          string
	Converted   bool
	Path        string
	Dst         D
	AfterFunc   func(D)
	OnEmptyFunc func()
}

func (cbg *convertibleBatchGetterItem[D]) After(f func(D)) ConvertibleBatchGetterItem[D] {
	cbg.AfterFunc = f
	return cbg
}

func (cbg *convertibleBatchGetterItem[D]) EmitAfter(d D) {
	if cbg.AfterFunc != nil {
		cbg.AfterFunc(d)
	}
}

func (cbg *convertibleBatchGetterItem[D]) RemoveAfter() ConvertibleBatchGetterItem[D] {
	cbg.AfterFunc = nil
	return cbg
}

func (cbg *convertibleBatchGetterItem[D]) OnEmpty(f func()) ConvertibleBatchGetterItem[D] {
	cbg.OnEmptyFunc = f
	return cbg
}

func (cbg *convertibleBatchGetterItem[D]) EmitEmpty() {
	if cbg.OnEmptyFunc != nil {
		cbg.OnEmptyFunc()
	}
}

func (cbg *convertibleBatchGetterItem[D]) RemoveOnEmpty() ConvertibleBatchGetterItem[D] {
	cbg.OnEmptyFunc = nil
	return cbg
}

type convertibleBatchGetter[S, D any] struct {
	BG      BatchGetter
	Items   sliceutil.Slice[*convertibleBatchGetterItem[*D]]
	DstMap  map[string]*D
	GetDoc  FuncGetDoc
	GetID   FuncGetID[D]
	Convert FuncConvert[S, D]
}

func NewConvertibleBatchGetter[S, D any](
	bg BatchGetter,
	getDoc FuncGetDoc,
	getID FuncGetID[D],
	convert FuncConvert[S, D],
) ConvertibleBatchGetter[S, D] {
	cbg := &convertibleBatchGetter[S, D]{
		BG:      bg,
		DstMap:  map[string]*D{},
		Items:   []*convertibleBatchGetterItem[*D]{},
		GetDoc:  getDoc,
		GetID:   getID,
		Convert: convert,
	}

	bg.OnCommit(func() {
		cbg.convertAll()
	})

	bg.OnEnd(func() {
		cbg.convertAll()
	})

	return cbg
}

func (cbg *convertibleBatchGetter[S, D]) convertAll() {
	items := cbg.Items.Filter(func(item *convertibleBatchGetterItem[*D]) bool {
		return !item.Converted && cbg.BG.IsCommittedItem(item.Path)
	})

	srcs := make([]any, len(items))

	// 変換処理
	for i, item := range items {
		src := cbg.BG.Get(item.Path)
		srcs[i] = src
		if src != nil {
			converted := cbg.Convert(src.(*S))
			if converted != nil {
				*(item.Dst) = *converted
			} else {
				srcs[i] = nil
			}
		}
	}

	// 変換が完了してからイベントを発火する
	for i, item := range items {
		src := srcs[i]
		if src == nil {
			item.EmitEmpty()
		} else {
			item.EmitAfter(item.Dst)
		}
		item.RemoveAfter()
		item.RemoveOnEmpty()
		cbg.DstMap[item.Path] = item.Dst
		item.Converted = true
	}
}

func (cbg *convertibleBatchGetter[S, D]) Add(ids ...string) ConvertibleBatchGetterItem[*D] {
	return cbg.SetWithID(new(D), ids...)
}

func (cbg *convertibleBatchGetter[S, D]) Set(d *D, ids ...string) ConvertibleBatchGetterItem[*D] {
	id := ""
	if d != nil {
		id = cbg.GetID(d)
	}
	return cbg.SetWithID(d, append(ids, id)...)
}

func (cbg *convertibleBatchGetter[S, D]) SetWithID(d *D, ids ...string) ConvertibleBatchGetterItem[*D] {
	docRef := cbg.GetDoc(ids...)
	data := new(S)
	cbg.BG.Add(docRef, data)

	item := &convertibleBatchGetterItem[*D]{
		ID:   docRef.ID,
		Path: docRef.Path,
		Dst:  d,
	}
	cbg.Items = append(cbg.Items, item)
	return item
}

func (cbg *convertibleBatchGetter[S, D]) Delete(ids ...string) {
	docRef := cbg.GetDoc(ids...)
	cbg.BG.Delete(docRef)
}

func (cbg *convertibleBatchGetter[S, D]) GetMap() map[string]*D {
	return maps.Clone(cbg.DstMap)
}

func (cbg *convertibleBatchGetter[S, D]) Get(ids ...string) *D {
	docRef := cbg.GetDoc(ids...)
	if dst, ok := cbg.DstMap[docRef.Path]; ok {
		return dst
	}
	return nil
}

func (cbg *convertibleBatchGetter[S, D]) Commit(ctx context.Context) error {
	return cbg.BG.Commit(ctx)
}
