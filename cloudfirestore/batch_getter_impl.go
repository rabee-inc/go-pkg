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
