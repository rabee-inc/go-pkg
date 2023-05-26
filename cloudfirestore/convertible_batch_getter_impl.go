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

func (a *convertibleBatchGetterItem[D]) After(f func(D)) ConvertibleBatchGetterItem[D] {
	a.AfterFunc = f
	return a
}

func (a *convertibleBatchGetterItem[D]) EmitAfter(d D) {
	if a.AfterFunc != nil {
		a.AfterFunc(d)
	}
}

func (a *convertibleBatchGetterItem[D]) RemoveAfter() ConvertibleBatchGetterItem[D] {
	a.AfterFunc = nil
	return a
}

func (a *convertibleBatchGetterItem[D]) OnEmpty(f func()) ConvertibleBatchGetterItem[D] {
	a.OnEmptyFunc = f
	return a
}

func (a *convertibleBatchGetterItem[D]) EmitEmpty() {
	if a.OnEmptyFunc != nil {
		a.OnEmptyFunc()
	}
}

func (a *convertibleBatchGetterItem[D]) RemoveOnEmpty() ConvertibleBatchGetterItem[D] {
	a.OnEmptyFunc = nil
	return a
}

type convertibleBatchGetter[S, D any] struct {
	bg      BatchGetter
	items   sliceutil.Slice[*convertibleBatchGetterItem[*D]]
	dstMap  map[string]*D
	getDoc  FuncGetDoc
	getID   FuncGetID[D]
	convert FuncConvert[S, D]
}

func NewConvertibleBatchGetter[S, D any](
	bg BatchGetter,
	getDoc FuncGetDoc,
	getID FuncGetID[D],
	convert FuncConvert[S, D],
) ConvertibleBatchGetter[S, D] {
	cbg := &convertibleBatchGetter[S, D]{
		bg:      bg,
		dstMap:  map[string]*D{},
		items:   []*convertibleBatchGetterItem[*D]{},
		getDoc:  getDoc,
		getID:   getID,
		convert: convert,
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
	items := cbg.items.Filter(func(item *convertibleBatchGetterItem[*D]) bool {
		return !item.Converted && cbg.bg.IsCommittedItem(item.Path)
	})

	srcs := make([]any, len(items))

	// 変換処理
	for i, item := range items {
		src := cbg.bg.Get(item.Path)
		srcs[i] = src
		if src != nil {
			*(item.Dst) = *(cbg.convert(src.(*S)))
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
		cbg.dstMap[item.Path] = item.Dst
		item.Converted = true
	}
}

func (cbg *convertibleBatchGetter[S, D]) Add(ids ...string) ConvertibleBatchGetterItem[*D] {
	return cbg.SetWithID(new(D), ids...)
}

func (cbg *convertibleBatchGetter[S, D]) Set(d *D, ids ...string) ConvertibleBatchGetterItem[*D] {
	id := ""
	if d != nil {
		id = cbg.getID(d)
	}
	return cbg.SetWithID(d, append(ids, id)...)
}

func (cbg *convertibleBatchGetter[S, D]) SetWithID(d *D, ids ...string) ConvertibleBatchGetterItem[*D] {
	docRef := cbg.getDoc(ids...)
	data := new(S)
	cbg.bg.Add(docRef, data)

	item := &convertibleBatchGetterItem[*D]{
		ID:   docRef.ID,
		Path: docRef.Path,
		Dst:  d,
	}
	cbg.items = append(cbg.items, item)
	return item
}

func (cbg *convertibleBatchGetter[S, D]) Delete(ids ...string) {
	docRef := cbg.getDoc(ids...)
	cbg.bg.Delete(docRef)
}

func (cbg *convertibleBatchGetter[S, D]) GetMap() map[string]*D {
	return maps.Clone(cbg.dstMap)
}

func (cbg *convertibleBatchGetter[S, D]) Get(ids ...string) *D {
	docRef := cbg.getDoc(ids...)
	if dst, ok := cbg.dstMap[docRef.Path]; ok {
		return dst
	}
	return nil
}

func (cbg *convertibleBatchGetter[S, D]) Commit(ctx context.Context) error {
	return cbg.bg.Commit(ctx)
}
