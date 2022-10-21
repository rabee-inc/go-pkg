package cloudfirestore

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/sliceutil"
	"github.com/rabee-inc/go-pkg/stringutil"
)

var reInvalidDocID = regexp.MustCompile(`\.|/`)

// ValidateDocumentID ... 正常な DocumentID かチェック
func ValidateDocumentID(str string) bool {
	return !reInvalidDocID.MatchString(str)
}

// ValidateCollectionRef ... 正常な Path かチェック
func ValidateCollectionRef(colRef *firestore.CollectionRef) bool {
	var docRef *firestore.DocumentRef
	for colRef != nil || docRef != nil {
		if colRef != nil {
			if !ValidateDocumentID(colRef.ID) {
				return false
			}
			docRef = colRef.Parent
			colRef = nil
		} else {
			if !ValidateDocumentID(docRef.ID) {
				return false
			}
			colRef = docRef.Parent
			docRef = nil
		}
	}
	return true
}

// ValidateDocumentRef ... 正常な Path かチェック
func ValidateDocumentRef(docRef *firestore.DocumentRef) bool {
	if !ValidateDocumentID(docRef.ID) {
		return false
	}
	return ValidateCollectionRef(docRef.Parent)
}

// GenerateDocumentRef ... ドキュメント参照を作成する
func GenerateDocumentRef(cFirestore *firestore.Client, docRefs []*DocRef) *firestore.DocumentRef {
	var dst *firestore.DocumentRef
	for i, docRef := range docRefs {
		if i == 0 {
			dst = cFirestore.Collection(docRef.CollectionName).Doc(docRef.DocID)
		} else {
			dst = dst.Collection(docRef.CollectionName).Doc(docRef.DocID)
		}
	}
	return dst
}

func RunTransaction(ctx context.Context, cFirestore *firestore.Client, fn func(ctx context.Context) error, opts ...firestore.TransactionOption) error {
	return cFirestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ctx = setContextTransaction(ctx, tx)
		return fn(ctx)
	}, opts...)
}

func RunWriteBatch(ctx context.Context, cFirestore *firestore.Client) context.Context {
	bt := cFirestore.Batch()
	return setContextWriteBatch(ctx, bt)
}

func CommitWriteBatch(ctx context.Context) (context.Context, error) {
	if bt := getContextWriteBatch(ctx); bt != nil {
		ctx = setContextWriteBatch(ctx, nil)
		if _, err := bt.Commit(ctx); err != nil {
			return ctx, err
		}
	} else {
		err := log.Errore(ctx, "no running write batch")
		return ctx, err
	}
	return ctx, nil
}

// Get ... 単体取得する(tx対応)
func Get(ctx context.Context, docRef *firestore.DocumentRef, dst interface{}) (bool, error) {
	if docRef == nil || docRef.ID == "" || !ValidateDocumentID(docRef.ID) {
		return false, nil
	}
	var dsnp *firestore.DocumentSnapshot
	var err error
	if tx := getContextTransaction(ctx); tx != nil {
		dsnp, err = tx.Get(docRef)
	} else {
		dsnp, err = docRef.Get(ctx)
	}
	if dsnp != nil && !dsnp.Exists() {
		return false, nil
	}
	if err != nil {
		log.Warning(ctx, err)
		return false, err
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		log.Error(ctx, err)
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	setEmptyBySlice(dst)
	return true, nil
}

// GetMulti ... 複数取得する(tx対応)
func GetMulti(ctx context.Context, cFirestore *firestore.Client, docRefs []*firestore.DocumentRef, dsts interface{}) error {
	docRefs = sliceutil.StreamOf(docRefs).
		Filter(func(docRef *firestore.DocumentRef) bool {
			return docRef != nil && docRef.ID != "" && ValidateDocumentID(docRef.ID)
		}).
		Out().([]*firestore.DocumentRef)
	if len(docRefs) == 0 {
		return nil
	}
	var dsnps []*firestore.DocumentSnapshot
	var err error
	if tx := getContextTransaction(ctx); tx != nil {
		dsnps, err = tx.GetAll(docRefs)
	} else {
		dsnps, err = cFirestore.GetAll(ctx, docRefs)
	}
	if err != nil {
		log.Warning(ctx, err)
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for _, dsnp := range dsnps {
		if !dsnp.Exists() {
			continue
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		setEmptyBySlices(rrv, rrt)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// GetByQuery ... クエリで単体取得する(tx対応)
func GetByQuery(ctx context.Context, query firestore.Query, dst interface{}) (bool, error) {
	query = query.Limit(1)
	var it *firestore.DocumentIterator
	if tx := getContextTransaction(ctx); tx != nil {
		it = tx.Documents(query)
	} else {
		it = query.Documents(ctx)
	}
	defer it.Stop()
	dsnp, err := it.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		log.Warning(ctx, err)
		return false, err
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		log.Error(ctx, err)
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	setEmptyBySlice(dst)
	return true, nil
}

// ListByQuery ... クエリで複数取得する(tx対応)
func ListByQuery(ctx context.Context, query firestore.Query, dsts interface{}) error {
	var it *firestore.DocumentIterator
	if tx := getContextTransaction(ctx); tx != nil {
		it = tx.Documents(query)
	} else {
		it = query.Documents(ctx)
	}
	defer it.Stop()
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for {
		dsnp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		setEmptyBySlices(rrv, rrt)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// ListByQueryCursor ... クエリで複数取得する（ページング）
func ListByQueryCursor(ctx context.Context, query firestore.Query, limit int, cursor *firestore.DocumentSnapshot, dsts interface{}) (*firestore.DocumentSnapshot, error) {
	if cursor != nil {
		query = query.StartAfter(cursor)
	}
	var it *firestore.DocumentIterator
	query = query.Limit(limit)
	if tx := getContextTransaction(ctx); tx != nil {
		it = tx.Documents(query)
	} else {
		it = query.Documents(ctx)
	}
	defer it.Stop()
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	var lastDsnp *firestore.DocumentSnapshot
	for {
		dsnp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Warning(ctx, err)
			return nil, err
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(v)
		if err != nil {
			log.Error(ctx, err)
			return nil, err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		setEmptyBySlices(rrv, rrt)
		rv.Set(reflect.Append(rv, rrv))
		lastDsnp = dsnp
	}
	if rv.Len() == limit {
		return lastDsnp, nil
	}
	return nil, nil
}

// Create ... 作成する(tx, bt対応)
func Create(ctx context.Context, colRef *firestore.CollectionRef, src interface{}) error {
	// 不正なIDがないかチェック
	if !ValidateCollectionRef(colRef) {
		return errors.New("Invalid Collection Path: " + colRef.Path)
	}
	setEmptyBySlice(src)
	var docRef *firestore.DocumentRef
	if tx := getContextTransaction(ctx); tx != nil {
		id := stringutil.UniqueID()
		docRef = colRef.Doc(id)
		err := tx.Create(docRef, src)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	} else if bt := getContextWriteBatch(ctx); bt != nil {
		id := stringutil.UniqueID()
		docRef = colRef.Doc(id)
		bt.Create(docRef, src)
	} else {
		var err error
		docRef, _, err = colRef.Add(ctx, src)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	}
	setDocByDst(src, docRef)
	return nil
}

// Update ... 更新する(tx, bt対応)
func Update(ctx context.Context, docRef *firestore.DocumentRef, kv map[string]interface{}) error {
	srcs := []firestore.Update{}
	for k, v := range kv {
		src := firestore.Update{Path: k, Value: v}
		srcs = append(srcs, src)
	}
	if tx := getContextTransaction(ctx); tx != nil {
		err := tx.Update(docRef, srcs)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	} else if bt := getContextWriteBatch(ctx); bt != nil {
		_ = bt.Update(docRef, srcs)
	} else {
		_, err := docRef.Update(ctx, srcs)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	}
	return nil
}

// Set ... 上書きする(tx, bt対応)
func Set(ctx context.Context, docRef *firestore.DocumentRef, src interface{}) error {
	// 不正なIDがないかチェック
	if !ValidateDocumentRef(docRef) {
		return errors.New("Invalid Document Path: " + docRef.Path)
	}
	setEmptyBySlice(src)
	if tx := getContextTransaction(ctx); tx != nil {
		err := tx.Set(docRef, src)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	} else if bt := getContextWriteBatch(ctx); bt != nil {
		_ = bt.Set(docRef, src)
	} else {
		_, err := docRef.Set(ctx, src)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	}
	setDocByDst(src, docRef)
	return nil
}

// Delete ... 削除する(tx, bt対応)
func Delete(ctx context.Context, docRef *firestore.DocumentRef) error {
	if tx := getContextTransaction(ctx); tx != nil {
		err := tx.Delete(docRef)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	} else if bt := getContextWriteBatch(ctx); bt != nil {
		_ = bt.Delete(docRef)
	} else {
		_, err := docRef.Delete(ctx)
		if err != nil {
			log.Warning(ctx, err)
			return err
		}
	}
	return nil
}

// AddStartWith ... 前方一致クエリを追加する
func AddStartWith(q firestore.Query, key string, word string) firestore.Query {
	return q.OrderBy(key, firestore.Asc).
		StartAt(word).
		EndAt(fmt.Sprintf("%s\uf8ff", word))
}
