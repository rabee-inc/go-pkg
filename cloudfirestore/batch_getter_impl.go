package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/maputil"
	"github.com/rabee-inc/go-pkg/sliceutil"
	"github.com/rabee-inc/go-pkg/util"
)

type batchGetterItem struct {
	docRef    *firestore.DocumentRef
	committed bool
	data      any
}

type batchGetter struct {
	docMap             map[string]*batchGetterItem
	commitEventEmitter util.EventEmitter
	endEventEmitter    util.EventEmitter
	cFirestore         *firestore.Client
}

func NewBatchGetter(cFirestore *firestore.Client) BatchGetter {
	return &batchGetter{
		docMap:             map[string]*batchGetterItem{},
		commitEventEmitter: util.NewEventEmitter(),
		endEventEmitter:    util.NewEventEmitter(),
		cFirestore:         cFirestore,
	}
}

func (bg *batchGetter) Add(docRef *firestore.DocumentRef, dst any) {
	if docRef == nil || docRef.ID == "" || !ValidateDocumentID(docRef.ID) {
		return
	}
	// 既に登録済みの場合は何もしない
	if _, ok := bg.docMap[docRef.Path]; ok {
		return
	}
	bg.docMap[docRef.Path] = &batchGetterItem{
		docRef: docRef,
		data:   dst,
	}
}

func (bg *batchGetter) Delete(docRef *firestore.DocumentRef) {
	delete(bg.docMap, docRef.Path)
}

func (bg *batchGetter) isAllCommitted(ctx context.Context) bool {
	for _, item := range bg.docMap {
		if !item.committed {
			return false
		}
	}
	return true
}

func (bg *batchGetter) commit(ctx context.Context) error {
	if len(bg.docMap) == 0 {
		return nil
	}
	docRefs := sliceutil.
		FilterMap(maputil.Values(bg.docMap), func(src *batchGetterItem) (bool, *firestore.DocumentRef) {
			return !src.committed, src.docRef
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
			d.committed = true
		} else {
			continue
		}
		err = dsnp.DataTo(dst)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
		SetDocByDst(dst, dsnp.Ref)
		SetEmptyBySlice(dst)
		SetEmptyByMap(dst)
	}
	return nil
}

func (bg *batchGetter) Commit(ctx context.Context) error {
	for !bg.isAllCommitted(ctx) {
		if err := bg.commit(ctx); err != nil {
			return err
		}
		bg.commitEventEmitter.Emit()
	}
	bg.commitEventEmitter.Clear()

	bg.endEventEmitter.Emit()
	bg.endEventEmitter.Clear()
	return nil
}

func (bg *batchGetter) IsCommittedItem(path string) bool {
	if d, ok := bg.docMap[path]; ok {
		return d.committed
	}
	return true
}

func (bg *batchGetter) Get(path string) any {
	if d, ok := bg.docMap[path]; ok {
		return d.data
	}
	return nil
}

func (bg *batchGetter) OnCommit(f func()) func() {
	return bg.commitEventEmitter.Add(f)
}

func (bg *batchGetter) OnEnd(f func()) func() {
	return bg.endEventEmitter.Add(f)
}
