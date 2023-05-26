package cloudfirestore

import "context"

type typedBatchGetter[T any] struct {
	bg     BatchGetter
	docMap map[string]string
	getDoc FuncGetDoc
}

func NewTypedBatchGetter[T any](bg BatchGetter, getDoc FuncGetDoc) TypedBatchGetter[T] {
	return &typedBatchGetter[T]{
		bg:     bg,
		docMap: map[string]string{},
		getDoc: getDoc,
	}
}

type FuncGetModel[E, M any] func(e *E) (id string, m *M)

func GetModelMapByBatchGetter[E, M any](
	bg TypedBatchGetter[E],
	getModel FuncGetModel[E, M],
) map[string]*M {
	es := bg.GetMap()
	ms := map[string]*M{}
	for _, e := range es {
		id, m := getModel(e)
		if m != nil {
			ms[id] = m
		}
	}
	return ms
}

func (tbg *typedBatchGetter[T]) Add(ids ...string) {
	docRef := tbg.getDoc(ids...)
	tbg.docMap[docRef.Path] = docRef.ID
	data := new(T)
	tbg.bg.Add(docRef, data)
}

func (tbg *typedBatchGetter[T]) Delete(ids ...string) {
	docRef := tbg.getDoc(ids...)
	delete(tbg.docMap, docRef.Path)
	tbg.bg.Delete(docRef)
}

func (tbg *typedBatchGetter[T]) GetMap() map[string]*T {
	m := map[string]*T{}
	for k, id := range tbg.docMap {
		d := tbg.bg.Get(k)
		if d != nil {
			m[id] = d.(*T)
		}
	}
	return m
}

func (tbg *typedBatchGetter[T]) Get(ids ...string) *T {
	docRef := tbg.getDoc(ids...)
	data := tbg.bg.Get(docRef.Path)
	return data.(*T)
}

func (tbg *typedBatchGetter[T]) Commit(ctx context.Context) error {
	return tbg.bg.Commit(ctx)
}
