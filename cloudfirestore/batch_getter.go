package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/maputil"
	"github.com/rabee-inc/go-pkg/sliceutil"
)

type docRefAndData struct {
	docRef *firestore.DocumentRef
	data   any
}

type BatchGetter interface {
	Add(docRef *firestore.DocumentRef, dst any)
	Delete(docRef *firestore.DocumentRef)
	Commit(ctx context.Context) error
	Get(path string) any
}

type batchGetter struct {
	docMap     map[string]*docRefAndData
	cFirestore *firestore.Client
}

func NewBatchGetter(cFirestore *firestore.Client) BatchGetter {
	return &batchGetter{
		docMap:     map[string]*docRefAndData{},
		cFirestore: cFirestore,
	}
}

func (bg *batchGetter) Add(docRef *firestore.DocumentRef, dst any) {
	if docRef == nil || docRef.ID == "" || !ValidateDocumentID(docRef.ID) {
		return
	}
	bg.docMap[docRef.Path] = &docRefAndData{
		docRef: docRef,
		data:   dst,
	}
}

func (bg *batchGetter) Delete(docRef *firestore.DocumentRef) {
	delete(bg.docMap, docRef.Path)
}

func (bg *batchGetter) Commit(ctx context.Context) error {
	if len(bg.docMap) == 0 {
		return nil
	}
	docRefs := sliceutil.
		Map(maputil.Values(bg.docMap), func(src *docRefAndData) *firestore.DocumentRef {
			return src.docRef
		})
	var dsnps []*firestore.DocumentSnapshot
	var err error

	if tx := getContextTransaction(ctx); tx != nil {
		dsnps, err = tx.GetAll(docRefs)
	} else {
		dsnps, err = bg.cFirestore.GetAll(ctx, docRefs)
	}
	if err != nil {
		log.Warning(ctx, err)
		return err
	}
	for _, dsnp := range dsnps {
		if !dsnp.Exists() {
			delete(bg.docMap, dsnp.Ref.Path)
			continue
		}
		var dst any
		if d, ok := bg.docMap[dsnp.Ref.Path]; ok {
			dst = d.data
		} else {
			continue
		}
		err = dsnp.DataTo(dst)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
		setDocByDst(dst, dsnp.Ref)
		setEmptyBySlice(dst)
	}
	return nil
}

func (bg *batchGetter) Get(path string) any {
	if d, ok := bg.docMap[path]; ok {
		return d.data
	}
	return nil
}

type FuncGetDoc func(ids ...string) *firestore.DocumentRef

type TypedBatchGetter[T any] interface {
	Add(ids ...string)
	Delete(ids ...string)
	GetMap() map[string]*T
	Get(ids ...string) *T
	Commit(ctx context.Context) error
}

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

type FuncBindBatchGetter[T any] func(bg BatchGetter) TypedBatchGetter[T]

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
